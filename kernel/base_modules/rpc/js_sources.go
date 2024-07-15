package kmodulerpc

const testJsProxySource = `(funct, proxyobject, ...parms) => {
	// Consolen Objekt
	console = { log: proxyobject.proxyShieldConsoleLog, error: proxyobject.proxyShieldErrorLog };

	// Timer funktionen
	clearInterval = () => proxyobject.clearInterval();
	clearTimeout = () => proxyobject.clearTimeout();
	setInterval = () => proxyobject.setInterval();
	setTimeout = () => proxyobject.setTimeout();

	// RPC Funktionen
	Resolve = (...parms) =>  proxyobject.resolve(...parms);

	// Testfunktionen für RPC Aufrufe
	Wait = (time) => proxyobject.wait(time);

	// Promise Proxy
	Promise = class CustodiaJsPromise extends Promise {
		constructor(executor) {
			const {resolveProxy, rejectProxy} = proxyobject.newPromise();
			const wrappedExecutor = (resolve, reject) => {
				executor(
					(value) => {
						resolveProxy();
						resolve(value);
					},
					(reason) => {
						rejectProxy();
						reject(reason);
					}
				);
			};
			super(wrappedExecutor);
		}
	}
	return funct(...parms);
}`
