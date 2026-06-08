export namespace gui {
	
	export class EncryptPayload {
	    inputPath: string;
	    outputPath: string;
	    readerPwd: string;
	    ownerPwd: string;
	    ownerOnly: boolean;
	    perms: string[];
	
	    static createFrom(source: any = {}) {
	        return new EncryptPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.outputPath = source["outputPath"];
	        this.readerPwd = source["readerPwd"];
	        this.ownerPwd = source["ownerPwd"];
	        this.ownerOnly = source["ownerOnly"];
	        this.perms = source["perms"];
	    }
	}

}

