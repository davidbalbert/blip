# Agent: Blip Parser (inspired by Carbon's parser)

Purpose
- Build a fast, robust, incremental-friendly parser that produces a postorder parse tree with strictly enforced structural invariants.
- Keep parse context shallow and defer most contextual/semantic checks to later stages.
- Maintain excellent error tolerance and recovery to support IDE use cases and refactoring tools.

Key takeaways from Carbon
- Data-oriented: store the parse tree as a flat, vectorized list in postorder for locality and fast downstream processing.
- Structural validity over grammatical validity: ensure every subtree has the right arity/bracketing shape even on invalid code; use has_error flags and explicit InvalidParse nodes for recovery.
- Minimal context in parse; push semantics and deep contextual validation to the check phase.
- Use explicit bracketing nodes and fixed child-count nodes to make the parse tree self-describing and easy to walk in postorder.

Scope
- This summary focuses on the parser's architecture, data model, state machine, node emission patterns (introducers, modifiers, optional clauses, operators), and error handling. It omits Carbon-specific grammar details and focuses on patterns Blip can reuse.

Interface and contracts
- Input: token buffer from the lexer (with kinds, lexeme spans, and positions). Optional precomputed bracket maps from lexing can be used for error recovery, not lookahead for normal parsing.
- Output: a flat parse tree in postorder; 1:1 node:token on valid input; extra nodes allowed for error recovery.
- Contract: every non-error subtree must have the exact bracketing or child-count described by its NodeKind. A Tree.Verify pass checks this.

Data model (Go)
- Token
  - Kind: enum
  - Span: {offset, length}, line/column for diagnostics
- NodeKind (enum)
  - Each kind defines either:
    - bracketing = true, or
    - childCount = fixed integer
- ParseNode
  - kind: NodeKind
  - token: token index (or none for synthetic nodes if you choose to allow them; Carbon avoids them except InvalidParse)
  - subtreeStart: index into nodes slice (points to first child in postorder window)
  - subtreeSize: uint32 (optional; try to avoid relying on it)
  - hasError: bool
- StateKind (enum)
  - The parser is a stack machine; states map to small handler functions
- ParseState
  - kind: StateKind
  - token: token index of introducer/anchor (when needed)
  - subtreeStart: node index snapshot captured when beginning a subtree
  - hasError: bool (bubble-up suppression)
- Context
  - tokens: []Token
  - pos: int (current token index)
  - nodes: []ParseNode (postorder storage)
  - stack: []ParseState
  - diagnostics: emitter

Core loop
- While stack not empty:
  - Pop next state, dispatch Handle<StateKind>(ctx)
  - Handlers:
    - Usually consume at least one token
    - Add nodes (AddLeafNode, AddNode)
    - Push follow-up states; the "next" state must be pushed last (LIFO)
- Expression parsing augments state entries with operator precedence info. Carbon stores this in the unified stack entry for efficiency.

Postorder tree as the storage format
- Emit nodes in postorder: children first, then their parent.
- Advantages:
  - Very cheap to produce in a single pass with minimal buffering.
  - Downstream stages (check) can process linearly with small local stacks; excellent cache locality.
- Bracketing nodes and child-count nodes are the essential contracts:
  - Bracketing node = opening token and first child in a subtree (example: "var", "fn", "{").
  - Fixed child-count node = e.g., infix "+" always has 2 children.

Bracketing and child-count contracts
- Each NodeKind declares either bracketing or childCount.
  - Bracketing nodes appear exactly once as the first child of their parent's child-slice, signaling subtree context to downstream consumers.
  - Fixed arity nodes enforce children count regardless of errors elsewhere; enables reliable traversal.
- Keep "=" or ";" as syntactic separators/roots when they clarify subtree boundaries useful for checking.

Error handling and invalid parses
- Strategy: produce a structurally valid tree that mirrors intended structure where possible; preserve information for IDE tooling.
- Use hasError flags:
  - Mark the directly erroneous node hasError = true.
  - Bubble up sparingly; only mark parents hasError if they can't proceed without a fully-checked child.
- Invalid nodes:
  - InvalidParse (leaf or synthetic parent to satisfy structure)
  - InvalidParseStart/InvalidParseSubtree patterns when you can't form the correct children
- Diagnostics:
  - Emit at first error site.
  - Use ReturnErrorOnState to suppress cascades and let parents skip follow-on unexpected tokens without re-diagnosing.
  - Always consume at least one token in an error path to avoid infinite loops.
  - SkipPastLikelyEnd (e.g., to next ";", closing brace, or outdent) to recover.

Parsing patterns to implement

1) Introducer pattern (dispatch on first token)
- Example: if, while, var, fn, etc.
- Implementation:
  - Switch on token kind at context.PositionKind()
    - On match:
      - Push state that will parse the rest (PushState)
      - Emit the introducer as a bracketing leaf (AddLeafNode)
    - Default:
      - Either delegate (e.g., to expression-statement), or emit diagnostic + InvalidParse, and consume one token
  - The closing token (often ";" or "}") becomes the parent node emitted in the finishing state, using subtreeStart captured when the introducer was seen.

2) Optional modifiers before introducer
- Example: virtual fn Foo();
- Problem: modifiers precede the introducer but you want the introducer to be the bracketing first child.
- Pattern:
  - Step 1: Capture subtreeStart = len(nodes)
  - Step 2: AddLeafNode(Placeholder, positionToken)
  - Step 3: Consume modifier tokens, emit modifier nodes as siblings
  - Step 4:
    - On introducer found: ReplacePlaceholderNode(subtreeStart, IntroducerKind, Consume())
    - On error: ReplacePlaceholderNode(subtreeStart, InvalidParseStart, pos, hasError=true), consume to likely end, then AddNode(InvalidParseSubtree, lastConsumed, subtreeStart, hasError=true)
  - Step 5: Push remaining states to parse the body/signature; ensure finishing state holds subtreeStart from step 1
  - Step 6: In finishing state, AddNode(DeclOrStmtRootKind, Consume(), subtreeStart, hasError)

3) Something required in context (required follower)
- Examples: a name after fn; a bracketed parameter list after a specific introducer.
- Approach:
  - The parent handler knows a child must appear next (e.g., identifier).
  - If absent, emit diagnostic; often emit InvalidParse or set hasError and recover.
  - This pattern is a small specialization: consume required token, or error and recover with a structurally valid substitute.

4) Optional clause (Case 1: introducer acts as parent of the clause)
- Example: function return type "-> Type"
- Emit flow:
  - If "->" present:
    - Push finishing state that will add ReturnType node using the saved "->" token
    - Consume "->" (ConsumeAndDiscard to delay node)
    - Push expression parse state for the type
  - Finisher:
    - AddNode(ReturnType, saved "->", subtreeStart, hasErrorFromChild)
- Without "->", don't emit the node.

5) Optional clause (Case 2: parent is required token after the optional clause, with distinct parent kinds)
- Example: impl [Type]? as Interface;
- Emit flow:
  - If "as" is next first:
    - AddLeafNode(DefaultSelfImplAs, Consume("as"))
  - Else:
    - Push state ImplBeforeAs, then parse the optional type expression
  - ImplBeforeAs finisher:
    - If "as" observed: AddNode(TypeImplAs, "as", subtreeStart includes type expr)
    - Else: if child had error, suppress new diag; otherwise, emit "expected as" diagnostic and ReturnErrorOnState
- This gives two different node kinds depending on option presence; it cleanly encodes the grammar in the tree.

6) Optional sibling (deprecated in Carbon docs)
- Carbon notes this pattern was changed; prefer Case 1 and Case 2 above to avoid fragile sibling dependencies.

7) Operators (precedence/associativity)
- Store operator precedence/associativity info inside the parse stack entries used by expression parsing.
- Use a shunting-yard-like or Pratt-style approach that emits nodes in postorder:
  - Operands are emitted first; operator nodes emitted when precedence dictates.
  - Infix nodes have childCount=2; prefix/postfix as needed.
- Keep expression parsing self-contained; no global lookahead beyond small fixed peeks.

Restrictive vs permissive parsing
- Favor permissive parsing:
  - Avoid arbitrary lookahead; 1–2 token lookahead is ok.
  - Avoid building complex context; keep parse fast/light and push to check.
  - Don't duplicate diagnostics in parse and check; centralize in check when semantics help.
- Parse still distinguishes syntax forms that would be hard to disambiguate later (e.g., places where choosing a form early changes the shape of the tree).

APIs (Go, suggested)
- Token navigation
  - PositionIs(kind), PositionKind(), Consume(), ConsumeIf(kind), ConsumeAndDiscard()
- Node emission
  - AddLeafNode(kind, tokenIndex [, hasError])
  - AddNode(kind, tokenIndex, subtreeStart [, hasError])
  - ReplacePlaceholderNode(index, kind, tokenIndex [, hasError])
- State management
  - PushState(kind [, tokenIndex=subtree anchor])
  - PopState() -> ParseState
  - PushStateForExpr(group) // precedence group for expressions
  - ReturnErrorOnState() // mark current top state as hasError
- Recovery
  - SkipPastLikelyEnd(semiOrDedentOrCloseBraceSet)

Verification
- Run Tree.Verify after parsing:
  - For each parent node:
    - If bracketing: must be first child; appears exactly once among siblings
    - If fixed childCount: exact number of children present
  - InvalidParse, InvalidParseStart/Subtree allowed where specified
- Fail fast in developer builds; convert to diagnostics in production if desired.

How this maps to Blip's incremental compiler
- Flat postorder nodes enable:
  - Fast linear passes in check and later phases
  - Cheap subtree invalidation/reparse: track token->node ranges; when tokens change, reparse covering ranges and splice nodes slice
- Stable indices:
  - Use integer indices into tokens and nodes; avoid pointers
  - For incremental edits, maintain a mapping of token spans to node ranges to bound reparses
- Data locality:
  - Single backing slices for tokens and nodes; reserve capacity to avoid realloc during normal builds

Implementation roadmap for Blip
1) Foundations
   - Define TokenKind, NodeKind, StateKind enums
   - Implement Context with nodes slice, stack, diagnostics
   - Implement the basic API methods (AddLeafNode, AddNode, ReplacePlaceholderNode, Push/PopState, Consume*, Verify)

2) Skeleton grammar
   - FileStart/FileEnd, block scopes with "{ … }"
   - Statement loop using Introducer pattern
   - Expression parser with precedence (infix, prefix, grouping)

3) Structural contracts
   - Mark bracketing kinds and childCount kinds in NodeKind metadata
   - Implement Verify enforcing contracts

4) Optional patterns
   - Modifiers-before-introducer (placeholder shuffling)
   - Optional clause Case 1 (fn return type)
   - Optional clause Case 2 (impl as-style with two parent kinds)

5) Error handling and recovery
   - Implement ReturnErrorOnState and hasError propagation
   - Implement SkipPastLikelyEnd and ensure single-token consumption on error paths
   - Add InvalidParse, InvalidParseStart/Subtree nodes and integrate into Verify

6) Integration with check
   - Document how check walks postorder:
     - On bracketing node: set context for the following children
     - On fixed arity nodes: pop children from a local stack
   - Ensure parse never relies on semantic info; only local disambiguation

7) Incremental support
   - Track node spans per token range
   - Implement minimal splice reparse for edited token ranges
   - Fuzz and property-test for structural validity under random edits

Testing strategy
- Golden tests comparing input code to flattened postorder node sequences
- Error recovery tests: deliberately malformed inputs; assert structural validity and minimal duplicate diagnostics
- Contract tests: Verify fails on constructed bad trees; passes on all emitted trees
- Operator precedence associativity tests
- Modifier/introducer shuffling tests

Performance notes
- Avoid allocations in handlers; pre-size nodes slice when possible
- Keep State and Node small, POD-like; prefer indices over pointers or strings
- Do not compute subtreeSize unless needed; if kept, compute once, use sparingly

Design trade-offs
- Postorder is optimal for downstream linear passes but requires careful thinking about "introducer now, parent later" emission patterns (handled via saved token + subtreeStart).
- Bracketing nodes add a node per subtree but pay off by making check simpler and faster.
- Permissive parse shifts some effort to check but avoids complex lookahead and improves IDE friendliness.

Summary
- Adopt Carbon's data-oriented parse: flat postorder tree, strict structure contracts via bracketing and fixed arity, and a small state-machine parser with explicit handlers.
- Keep parse permissive, context-light, and error-tolerant; push semantic checks to later phases.
- Implement the specific node emission patterns (introducer, modifiers-before-introducer, optional clauses Case 1/2, operator precedence) to achieve a robust, incremental-friendly Blip parser.

Reference: https://docs.carbon-lang.dev/toolchain/docs/parse.html
