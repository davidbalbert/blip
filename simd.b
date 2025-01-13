// SIMD questions
// - Should have have custom SIMD types or just use arrays?
// - How much APLish array programming should we support (I think at least some!)
// - Tension between things that work well with SIMD vs GPUs (if we ever want to target shaders)
// - For vectors, should we support things like v.xy, v.yz, etc.?

// If we have some custom SIMD types, here are some ideas for names:

vNintM
vNfloatM

E.g.

v2int32
v3float64

Aliases for 
[2]int32
[3]int64

// But have all combinations of x, y, z, w (e.g. .xz) returns a v2. 

// Alt names

v3int32
vec3int32
int32v3
int32vec3
int32x3
v[3]int32
vec[3]int32
vec3[int32]
[vec3]int32
int32[vec3]
[int32]vec3

v3i32
vec3i32
i32vec3
// etc.


m34int32
m3x4int32
mat34int32
mat3x4int32
int32m34
int32m3x4
int32mat34
int32m3x4
int32x3x4
m[3,4]int32
mat[3,4]int32
m[3x4]int32
mat[3x4]int32

m34i32
m3x4int32
// etc.

// Or just use slicing operations. E.g. m[1,2] returns a v2.