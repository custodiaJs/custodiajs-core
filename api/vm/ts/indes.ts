// Es wird geprüft ob es sich um eine VNH1 Umgebung handelt
if (typeof vnh1 === undefined) throw new Error("not supported runtime");

declare namespace vnh1 {
    function com(message: string, ...data:any): any;
    let version:string;
}

// Die VM Importe werden imporiert
const vnh1VMRootSignal = vnh1.com("root/vm");
const vnh1VMModule = vnh1.com("root/modules");

// Console Imports
const vnh1ConsoleLog = vnh1.com("console/log");
const vnh1ConsoleInfo = vnh1.com("console/info");
const vnh1ConsoleError = vnh1.com("console/error");

// Share Function Imports
const vnha1ShareFunction = vnh1.com("root/sharefunction");

// S3 Imports
const vnh1S3Init = vnh1.com("s3/initobject");
const vnh1S3Upload = vnh1.com("s3/initobject");
const vnh1S3Download = vnh1.com("s3/initobject");
const vnh1S3Delete = vnh1.com("s3/initobject");

// HTTP Importe
const vnh1HttpClientGet = vnh1.com("s3/initobject");
const vnh1HttpClientPost = vnh1.com("s3/initobject");
const vnh1HttpClientPut = vnh1.com("s3/initobject");
const vnh1HttpClientDelete = vnh1.com("s3/initobject");

// Cache Importe
const vnh1CacheWrite = vnh1.com("s3/initobject");
const vnh1CacheRead = vnh1.com("s3/initobject");

// Share Function Exports
export function shareFunction(functionName: string, passedFunction:Function) {
    // Es wird geprüft ob die Sharing Funktion aktiv ist
    if (!vnh1VMModule("function_share")) throw new Error("function sharing is disabeld");

    // Die Funktion wird geteilt
    try {vnha1ShareFunction(functionName, passedFunction)}
    catch(e) {throw e;}
}

// Console Exports
export const console = {
    log : (...args: any[]): void => vnh1ConsoleLog(...args),
    info : (...args: any[]): void => vnh1ConsoleInfo(...args),
    error : (...args: any[]): void => vnh1ConsoleError(...args),
}

// S3 Exports
export interface S3Object {
    key: string;
    data: string;
    metadata: Record<string, string>; // Metadatenfelder als Schlüssel-Wert-Paare
}

export class S3Client {
    private registerId: number;

    constructor(bucketName: string) {
        // Es wird geprüft ob der S3 Dienst verfügbar ist
        if (!vnh1VMModule("s3")) throw new Error("s3 is disabeld");

        // Der Vorgang wird registriert
        var result:any;
        try{result=vnh1S3Init(bucketName);}
        catch(e) {}
    
        // Die ID wird zwischengespeichert
        this.registerId = result;
    }

    async uploadObject(key: string, data: string, metadata: Record<string, string>): Promise<void> {
        // Es wird geprüft ob der S3 Dienst verfügbar ist
        if (!vnh1VMModule("s3")) throw new Error("s3 is disabeld");
        await vnh1S3Upload(this.registerId, key, data, metadata)
    }

    async downloadObject(key: string, metadata: Record<string, string>): Promise<S3Object | null> {
        // Es wird geprüft ob der S3 Dienst verfügbar ist
        if (!vnh1VMModule("s3")) throw new Error("s3 is disabeld");
        return await vnh1S3Download(this.registerId, key, metadata)
    }

    async deleteObject(key: string, metadata: Record<string, string>): Promise<void> {
        // Es wird geprüft ob der S3 Dienst verfügbar ist
        if (!vnh1VMModule("s3")) throw new Error("s3 is disabeld");
        return await vnh1S3Delete(this.registerId, key, metadata)
    }
}

// HTTP Exporte
export interface HttpRequestOptions {
    headers?: Record<string, string>;
    body?: any;
    queryParams?: URLSearchParams | Record<string, string>;
}

export interface HttpResponse<T = any> {
    status: number;
    statusText: string;
    headers: Record<string, string>;
    data: T;
}

export interface HttpClient {
    get<T>(url: string, options?: HttpRequestOptions): Promise<HttpResponse<T>>;
    post<T>(url: string, options?: HttpRequestOptions): Promise<HttpResponse<T>>;
    put<T>(url: string, options?: HttpRequestOptions): Promise<HttpResponse<T>>;
    delete<T>(url: string, options?: HttpRequestOptions): Promise<HttpResponse<T>>;
    // Zusätzliche Methoden können hier hinzugefügt werden, z.B. PATCH
}

// Der VM wird Signalisiert dass der Vorgang erfolgreich durchgeführt wurde
vnh1VMRootSignal("100000");
