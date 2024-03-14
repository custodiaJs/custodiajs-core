if (typeof vnh1 === undefined) throw new Error("not supported runtime");

declare namespace vnh1 {
    function com(message: string, ...data:any): any;
    let version:string;
}

// Die VM Importe werden imporiert
const vnh1VMRootSignal = vnh1.com("root/vm");

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

// Share Function Exports
export function shareFunction(functionName: string, passedFunction:Function) {
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
        // Der Vorgang wird registriert
        var result:any;
        try{result=vnh1S3Init(bucketName);}
        catch(e) {}
    
        // Die ID wird zwischengespeichert
        this.registerId = result;
    }

    async uploadObject(key: string, data: string, metadata: Record<string, string>): Promise<void> {
        await vnh1S3Upload(this.registerId, key, data, metadata)
    }

    async downloadObject(key: string, metadata: Record<string, string>): Promise<S3Object | null> {
        return await vnh1S3Download(this.registerId, key, metadata)
    }

    async deleteObject(key: string, metadata: Record<string, string>): Promise<void> {
        return await vnh1S3Delete(this.registerId, key, metadata)
    }
}

// Der VM wird Signalisiert dass der Vorgang erfolgreich durchgeführt wurde
vnh1VMRootSignal("100000");
