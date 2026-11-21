// Starts the Go client. wasm_exec.js defines globalThis.Go.
import "./wasm_exec.js";

const go = new Go();
const { instance } = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject);
go.run(instance);
