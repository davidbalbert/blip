// One of our goals is "Perfect C interop". What does this mean?
//
// 1. Call any C function from Blip with minimal ceremony.
// 2. Call any Blip function from C. Possibly even functions that take non-bridgable types like
//    tagged unions.
//    - Automatic header generation.
// 3. Exist as a peer to C rather than a above it. The compiler should be able to generate .o files
//    which can be linked with existing .o files generated by C using Make, CMake, etc.
// 4. Easily convert unsafe C pointers (!*T) into the correct smart pointer type.
//    - Wrap foreign reference counted types in a #*T, including support for auto-niling weak pointers.
//    - Cast !*T to $*T with a custom free function, so that the value can be disposed of appropriately.
//    - Mark external types as NoCopy so Blip can enforce cleanup.
// 5. A text format for out of band annotations of C declarations so they get imported with smart
//    pointer types.
//    - Maybe also include regexes so that we can annotate a bunch of functions at once, e.g. to apply
//      the "Create Rule" from Core Foundation where all values returned by functions that have "Create"
//      or "Copy" in the name indicate that you own the returned object (it has a +1 refcount) and must
//      release it with CFRelease.
//    - The ability to annotate C types (e.g. the int returned by open(2) is no-copy). This might be infeasable:
//      it would probably require getting rid of mem.NoCopy (which is easy), but file descriptors can be copied –
//      e.g. write(2) does not consume an fd. Perhaps we should not attempt this. The user is always responsible
//      for making sure resources are deallocated. Otherwise, we'd have to allow some functions to borrow
//      non-pointer values.
// 6. Make C safer.
//    - Nil checks even for pointers returned from C (easy).
//    - Bounds checks or alias warnings even for pointers passed to C (hard).
// 8. First class support for "alien" memory objects.
//    - Memory mapped registers at known addresses
//    - Memory populated by the OS, e.g. Linux vDSO, the auxiliary vector, etc.
// 9. No need to use any of the above safety features. You should be able to just call C functions as
//    they are, with no Blip interop affordances. The code will be unsafe – you have to enforce
//    whatever invariants exist manually – but that's ok.
//
// For inspiration, see "Some Were Meant for C" - https://www.humprog.org/~stephen/research/papers/kell17some-preprint.pdf
