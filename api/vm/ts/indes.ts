// Es wird geprüft ob es sich um eine VNH1 Umgebung handelt
if (typeof vnh1 === undefined) throw new Error("not supported runtime");

// Die VNh1 Namespaces werden deklariert
declare namespace vnh1 {
    function com(message: string, ...data:any): any;
    let version:string;
}

// Das S3 Interface für die Bridge wird erezugt
interface S3IoBridge {
    uploadObject(key: string, data: any, metadata: Record<string, string>): Promise<void>
    downloadObject(key: string, metadata: Record<string, string>): Promise<S3Object | null> 
    deleteObject(key: string, metadata: Record<string, string>): Promise<void>
}

// Console Imports
const consoleModule = vnh1.com("console");

// Cache Importe
const cacheModule = vnh1.com("cache");

// Die VM Importe werden imporiert
const rootModule = vnh1.com("root");

// S3 Imports
const s3Module = vnh1.com("s3");

// Console Export
export const console = {
    log : function(...args:any) {
        consoleModule("log", ...args);
    },
    info : function(...args:any) {
        consoleModule("info", ...args);
    },
    error : function(...args:any) {
        consoleModule("error", ...args);
    }
}

// Cache Export
export const cache = {
    write: function(name:string, ...args:any) {
        cacheModule("write", name, ...args);
    },
    read: function(name:string):any {
        return cacheModule("read", name);
    },
}

// S3 Export
export interface S3Object {
    key: string;
    data: string;
    metadata: Record<string, string>; // Metadatenfelder als Schlüssel-Wert-Paare
}

// S3 Client Export
export class S3Client {
    private ioBridge: S3IoBridge;

    constructor(bucketNameOrUrl: string) {
        // Es wird geprüft ob der S3 Dienst verfügbar ist
        if (!rootModule("mavail", "s3")) throw new Error("s3 is disabeld");

        // Es wird geprüft ob die S3 Modul Funktionen bereitstehen
        if (s3Module === undefined) throw new Error("s3 is disabeld");

        // Der Vorgang wird registriert
        try{this.ioBridge = (s3Module("init", bucketNameOrUrl) as S3IoBridge);}
        catch(e) { throw e; }
    }

    async uploadObject(key: string, data: string | number | ArrayBuffer, metadata: Record<string, string>): Promise<void> {
        // Es wird geprüft ob der S3 Dienst verfügbar ist
        if (!rootModule("mavail", "s3")) throw new Error("s3 is disabeld");

        // Es wird geprüft ob die S3 Modul Funktionen bereitstehen
        if (s3Module === undefined) throw new Error("s3 is disabeld");

        await this.ioBridge.uploadObject(key, data, metadata);
    }

    async downloadObject(key: string, metadata: Record<string, string>): Promise<S3Object | null> {
        // Es wird geprüft ob der S3 Dienst verfügbar ist
        if (!rootModule("mavail", "s3")) throw new Error("s3 is disabeld");

        // Es wird geprüft ob die S3 Modul Funktionen bereitstehen
        if (s3Module === undefined) throw new Error("s3 is disabeld");

        return await this.ioBridge.downloadObject(key, metadata)
    }

    async deleteObject(key: string, metadata: Record<string, string>): Promise<void> {
        // Es wird geprüft ob der S3 Dienst verfügbar ist
        if (!rootModule("mavail", "s3")) throw new Error("s3 is disabeld");

        // Es wird geprüft ob die S3 Modul Funktionen bereitstehen
        if (s3Module === undefined) throw new Error("s3 is disabeld");

        return await this.ioBridge.deleteObject(key, metadata)
    }
}

// Share Function Export
export function localFunctionShare(functionName: string, passedFunction:Function) {
    // Es wird geprüft ob die Sharing Funktion aktiv ist
    if (!rootModule("mavail", "function_share")) throw new Error("function sharing is disabeld");

    // Die Funktion wird geteilt
    try {rootModule("fshare", "local", functionName, passedFunction)}
    catch(e) {throw e;}
}

// Share Function Export
export function publicFunctionShare(functionName: string, passedFunction:Function) {
    // Es wird geprüft ob die Sharing Funktion aktiv ist
    if (!rootModule("mavail", "function_share")) throw new Error("function sharing is disabeld");

    // Die Funktion wird geteilt
    try {rootModule("fshare", "local", functionName, passedFunction)}
    catch(e) {throw e;}
}

// Es wird der VM Signalisiert dass die Initalisierung der API Erfolgreich abgeschlossen wurde
if (!rootModule("finsh")) throw new Error("api initalization failed");

// Tests
const test = new S3Client("uri");
test.uploadObject("test", "data", {"arga":"value"});

localFunctionShare("test", () => {
    console.log("test");
});

cache.write("test", true)
const a = cache.read("test");
console.log(a)
