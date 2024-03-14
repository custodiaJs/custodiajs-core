"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.WebSocket = exports.console = void 0;
// Console
exports.console = {
    log: (...args) => vnh1.com("console/log", ...args),
    info: (...args) => vnh1.com("console/info", ...args),
    error: (...args) => vnh1.com("console/error", ...args),
};
class WebSocket {
    // Konstruktor
    constructor(url) {
        this.url = url;
        // Ereignishandler
        this.onopenHandler = null;
        this.onmessageHandler = null;
        this.onerrorHandler = null;
        this.oncloseHandler = null;
        exports.console.log(`WebSocket connection to '${url}' will be simulated.`);
    }
    // Methoden zum Setzen der Ereignishandler
    set onopen(handler) {
        this.onopenHandler = handler;
    }
    set onmessage(handler) {
        this.onmessageHandler = handler;
    }
    set onerror(handler) {
        this.onerrorHandler = handler;
    }
    set onclose(handler) {
        this.oncloseHandler = handler;
    }
    // Methode zum Simulieren des Sendens einer Nachricht
    send(data) {
        exports.console.log(`Sending message: ${data}`);
        // Simulieren Sie eine Antwort vom Server nach einer kurzen Verzögerung
        setTimeout(() => {
            var _a;
            (_a = this.onmessageHandler) === null || _a === void 0 ? void 0 : _a.call(this, { data: `Echo: ${data}` });
        }, 500);
    }
    // Methode zum Simulieren des Öffnens der Verbindung
    open() {
        var _a;
        exports.console.log(`Simulating open WebSocket connection to ${this.url}`);
        (_a = this.onopenHandler) === null || _a === void 0 ? void 0 : _a.call(this);
    }
    // Methode zum Simulieren des Schließens der Verbindung
    close() {
        var _a;
        exports.console.log(`Closing WebSocket connection to ${this.url}`);
        (_a = this.oncloseHandler) === null || _a === void 0 ? void 0 : _a.call(this, { code: 1000, reason: "Normal closure" });
    }
    // Methode zum Simulieren eines Fehlers
    error() {
        var _a;
        exports.console.log(`Simulating WebSocket error`);
        (_a = this.onerrorHandler) === null || _a === void 0 ? void 0 : _a.call(this, new Error("Simulated error"));
    }
}
exports.WebSocket = WebSocket;
exports.console.log("test", "abc");
