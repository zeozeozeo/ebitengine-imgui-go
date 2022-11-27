# Ebitengine rendering backend for the [pure Go imgui port](https://github.com/Splizard/imgui)

This uses a [fork](https://github.com/zeozeozeo/imgui) of the version by [Splizard](https://github.com/Splizard), which adds some minor changes to it.

# Benefits of not using CGo

-   Cross-compilation
-   WebAssembly
-   Pure Go :D

# Building the example for WebAssembly

1. Clone the repository
2. Navigate into the example directory: `cd internal/example` (on Linux)
3. Build the example:

    On Linux:

    ```
    env GOOS=js GOARCH=wasm go build -o example.wasm main.go
    ```

    On Windows Powershell:

    ```
    $Env:GOOS = 'js'
    $Env:GOARCH = 'wasm'
    go build -o yourgame.wasm github.com/yourname/yourgame
    Remove-Item Env:GOOS
    Remove-Item Env:GOARCH
    ```

4. Copy `wasm_exec.js` into the current directory:

    On Linux:

    ```
    cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
    ```

    On Windows Powershell:

    ```
    $goroot = go env GOROOT
    cp $goroot\misc\wasm\wasm_exec.js .
    ```

5. Create this HTML file

    ```html
    <!DOCTYPE html>
    <script src="wasm_exec.js"></script>
    <script>
        // Polyfill
        if (!WebAssembly.instantiateStreaming) {
            WebAssembly.instantiateStreaming = async (resp, importObject) => {
                const source = await (await resp).arrayBuffer();
                return await WebAssembly.instantiate(source, importObject);
            };
        }

        const go = new Go();
        WebAssembly.instantiateStreaming(
            fetch("example.wasm"),
            go.importObject
        ).then((result) => {
            go.run(result.instance);
        });
    </script>
    ```

6. Start a local HTTP server and open the page in your browser

If you want to embed the game into another page, use iframes (assuming that `main.html` is the name of the above HTML file):

```
<!DOCTYPE html>
<iframe src="main.html" width="640" height="480" allow="autoplay"></iframe>
```
