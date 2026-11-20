# Jani App — Go + C++ + Slint

> **Logic in Go, Design in Slint, Bridge in C++.**
> Build the bridge once, then never touch it again. After that, work only in Go and Slint.

---

## 1. What Is This Project?

A light, fast desktop app.

| Part       | Job                               | Language |
| :--------- | :-------------------------------- | :------- |
| **Go**     | All logic, data, callbacks        | Go       |
| **Slint**  | Only UI / Design                  | Slint    |
| **C++**    | The door between Go and Slint     | C++      |
| **Bridge** | Generic API, no hardcoded screens | C++ + Go |

**Size:** ~23 MB total

- Slint engine: ~21 MB
- Bridge: ~44 KB (less than 0.2%)
- So adding features to the bridge does not affect size.

**Rules:**

- No logic in UI.
- No UI in Go.
- No screen names hardcoded in C++. C++ only moves data back and forth.

---

## 2. Folder Structure

```text
slint-go/
├── bridge/
│   ├── bridge.cpp
│   ├── bridge.h
│   └── CMakeLists.txt
├── build/
│   ├── dev/              # Dev build output - libbridge.so
│   └── release/          # Release build output - libbridge.so
├── go/
│   ├── bridge.go         # Core - Init, Load, Run, SetString, OnCallback
│   ├── callbacks.go      # GoCallbackHandler - Slint to Go callback route
│   ├── go.mod
│   ├── logic.go          # Your logic - LogicDoWork
│   ├── main.go           # App entry point
│   └── types.go          # Arrays + Structs - SetStructArray etc.
├── ui/
│   ├── app.slint         # Slint UI - design only
│   └── avatar.png
├── build.log
└── README.md
```

---

## 3. Bridge API — What Is Available?

The bridge is **generic**. It does not know if the data is a `User` or a `Product`. It just passes data.

### 3.1 Core — `go/bridge.go`

| Go Function                           | What Goes to Slint              |
| :------------------------------------ | :------------------------------ |
| `SetString("my-data", "Hello")`       | `string`                        |
| `SetInt("count", 5)`                  | `int`                           |
| `SetBool("is-loading", true)`         | `bool`                          |
| `SetFloat("progress", 0.5)`           | `float`                         |
| `SetColor("bg-color", "#eb1b0c")`     | `color`                         |
| `SetImage("avatar", "ui/avatar.png")` | `image` - from local file       |
| `SetImageBytes("avatar", bytes)`      | `image` - from downloaded bytes |
| `GetString("my-data")`                | Get `string` from Slint         |
| `GetInt("count")`                     | Get `int` from Slint            |
| `GetBool("is-loading")`               | Get `bool` from Slint           |
| `OnCallback("do-work", func)`         | Slint button -> Go logic        |

### 3.2 Types — `go/types.go`

| Go Function                                    | What Goes to Slint           |
| :--------------------------------------------- | :--------------------------- |
| `SetStringArray("names", []string{...})`       | `[string]`                   |
| `SetIntArray("scores", []int{...})`            | `[int]`                      |
| `SetFloatArray("prices", []float64{...})`      | `[float]`                    |
| `SetStruct("user", struct)`                    | One `struct`                 |
| `SetStructArray("users", []*SlintStruct{...})` | `[Struct]` - list of structs |
| `SetEnum("mode", "dark")`                      | Sends `enum` as `string`     |

### 3.3 Struct — Flat Is Enough

> **Note:** Struct of Struct (struct inside struct) is not required. 90% of cases work with a flat struct. We will add `struct_set_struct` only when nested data is really needed.

```go
u := NewStruct()
u.SetString("name", "Abdullah")
u.SetInt("age", 19)
u.SetString("city", "Wah")
u.SetBool("is-vip", true)
b.SetStruct("user", u)
u.Free()
```

```slint
struct User {
    name: string,
    age: int,
    city: string,
    is-vip: bool,
}
in property <User> user;
Text { text: user.name; }
```

### 3.4 Array of Structs — List

```go
u1 := NewStruct()
u1.SetString("name", "Abdullah")
u1.SetInt("age", 19)

u2 := NewStruct()
u2.SetString("name", "Ali")
u2.SetInt("age", 20)

b.SetStructArray("users", []*SlintStruct{u1, u2})
u1.Free(); u2.Free()
```

```slint
struct User { name: string, age: int, }
in property <[User]> users;

for user in users: Rectangle {
    Text { text: user.name + " - " + user.age; }
}
```

### 3.5 Callback — Slint to Go

```go
b.OnCallback("do-work", func(input string) string {
    b.SetStructArray("users", users)
    return "Done"
})
```

```slint
callback do-work(string) -> string;
TouchArea {
    clicked => { root.my-data = root.do-work(input.text); }
}
```

> **Remember:** Register `OnCallback` **before** `LoadUI`. Set props **after** `LoadUI`.

---

## 4. Dev Setup — One Time (Arch Linux)

### 4.1 Packages

```bash
sudo pacman -S --needed base-devel cmake git go fontconfig libxkbcommon wayland libxcb rustup
rustup default stable
```

Check:

```bash
gcc --version
cmake --version
go version
cargo --version
```

### 4.2 Rust Flag — Required

Slint 1.8.0 Rust code fails with the new Rust compiler strict lint `dangerous_implicit_autorefs`. This lint is only inside Slint, not in your code.

Before every build in a new terminal:

```bash
export RUSTFLAGS="-A dangerous_implicit_autorefs"
```

To make it permanent, add this line to `~/.bashrc` or `~/.zshrc`:

```bash
export RUSTFLAGS="-A dangerous_implicit_autorefs"
```

If a new Slint lint error appears, add it the same way:

```bash
export RUSTFLAGS="-A dangerous_implicit_autorefs -A some_other_lint"
```

---

## 5. Dev Build — Daily Work

From project root (the folder with `ui`, `go`, `bridge`):

### Step 1 — Build C++ Bridge

```bash
export RUSTFLAGS="-A dangerous_implicit_autorefs"

cmake -S bridge -B build/dev -DCMAKE_BUILD_TYPE=Debug
cmake --build build/dev -j$(nproc) 2>&1 | tee build.log
```

- First time Slint downloads and compiles, it can take 5 to 10 minutes.
- Warnings are fine. Only lines starting with `error:` matter.
- On success you will see `[100%] Built target bridge`.

Check:

```bash
ls build/dev/libbridge.so
```

Check for errors:

```bash
grep -n "^error" build.log | head -5
```

### Step 2 — Run Go App

```bash
cd go
LD_LIBRARY_PATH=../build/dev go run .
```

- `LD_LIBRARY_PATH` tells the app where `libbridge.so` is. It is only needed at run time.
- Window will open. Press `Do Work`, type text to test.

### When to Rebuild?

| What Changed?                  | What to Do?                                  |
| :----------------------------- | :------------------------------------------- |
| Only `.slint` or Go code       | No `cmake` needed, just `go run`             |
| `bridge/*.cpp` or `bridge/*.h` | Run `cmake --build build/dev`, then `go run` |
| Build folder looks broken      | `rm -rf build/dev` and repeat Step 1         |

> **Tip:** Dev uses the Slint interpreter for fast reload. Release compiles Slint into C++. The bridge / API is the same in both, Go does not know the difference.

---

## 6. Release Build — Final App

In dev, `LD_LIBRARY_PATH` works, but in release the app should run from any folder. So we use `$ORIGIN` rpath, so the app finds `.so` files in its own folder.

From project root:

```bash
export RUSTFLAGS="-A dangerous_implicit_autorefs"

# 1. Create dist folder
mkdir -p dist

# 2. Bridge Release build
cmake -S bridge -B build/release -DCMAKE_BUILD_TYPE=Release
cmake --build build/release -j$(nproc)

# 3. Build Go app, $ORIGIN makes it look for .so in the same folder
cd go
go build -ldflags='-s -w -r $ORIGIN' -o ../dist/app .
cd ..
```

### Copy Files

```bash
cp build/release/libbridge.so dist/
cp build/release/_deps/slint-build/libslint_cpp.so dist/ 2>/dev/null || cp $(find build/release -name libslint_cpp.so | head -1) dist/
cp -r ui dist/
```

Now `dist/` will contain:

```text
dist/
├── app                # Final app
├── libbridge.so
├── libslint_cpp.so
└── ui/
    └── app.slint
```

### Run

```bash
cd dist && ./app
```

> **Note:** Copying `ui/` to `dist/` is required, because the `.slint` file is loaded at runtime.

If the app does not run, check missing libraries:

```bash
ldd dist/app
ldd dist/libbridge.so
```

If you see a `not found` line, check that library path.

---

## 7. Slint Design — Without Import

Importing from `std-widgets.slint` often gives this error:

```text
Slint error: Syntax error: expected Identifier
Slint error: No exported type called '' found in "std-widgets.slint"
```

Fix: Skip `std-widgets`, use Slint built-in components. They never fail.

| Before (Fails)  | Now (Built-in, Safe)               |
| :-------------- | :--------------------------------- |
| `VerticalBox`   | `VerticalLayout`                   |
| `HorizontalBox` | `HorizontalLayout`                 |
| `LineEdit`      | `TextInput`                        |
| `Button`        | `Rectangle` + `TouchArea` + `Text` |

```slint
export component App inherits Window {
    // No import needed
    VerticalLayout {
        TextInput { text: ""; }
        Rectangle {
            Text { text: "Do Work"; }
            TouchArea { clicked => { /* ... */ } }
        }
    }
}
```

---

## 8. Common Errors

| Error                                          | Cause                     | Fix                                                                 |
| :--------------------------------------------- | :------------------------ | :------------------------------------------------------------------ |
| `dangerous_implicit_autorefs` error            | RUSTFLAGS not set         | Run `export RUSTFLAGS="-A dangerous_implicit_autorefs"`, then build |
| `build/dev is not a directory`                 | Wrong build folder        | `rm -rf build/dev`, then configure again                            |
| `libbridge.so: cannot open shared object`      | LD_LIBRARY_PATH missing   | Run `LD_LIBRARY_PATH=../build/dev go run .`                         |
| `UI load fail` / Slint error                   | Syntax in `ui/app.slint`  | Look at `Slint error:` message in terminal                          |
| `cargo not found`                              | Rust not installed        | `rustup default stable`                                             |
| Window does not open, `ui/app.slint` not found | Running from wrong folder | Run `go run .` from `go/` folder                                    |
| `expected Identifier` / `std-widgets` error    | Import issue              | Use Section 7 no-import design                                      |

Warnings from Rust / Slint can be ignored. Only focus on `error`.

---

## 9. What Can Be Added Next?

Bridge is only 44 KB, there is plenty of room. These can be added later without size worries:

- [ ] File Drop — Single / Multiple / Folder
- [ ] System Tray
- [ ] Notification
- [ ] Clipboard
- [ ] Screenshot
- [ ] File Picker
- [ ] Video — Via OS player (not inside app)

Video codecs to keep: MP4 + H264 + AAC, MP4 + H265 + AAC, WebM + VP9. AV1 is dropped for now.

---

## 10. Golden Rules

1. **Build the bridge once, then never touch it.**
2. **Logic in Go, Design in Slint.** Do not mix them.
3. **No hardcoded screen names in C++.**
4. **Register callbacks before `LoadUI`, set props after `LoadUI`.**
5. **Use flat structs.** Use struct-of-struct only when data is truly nested.

---

_App — Built with Go + C++ + Slint. Fast, light, and fun._
