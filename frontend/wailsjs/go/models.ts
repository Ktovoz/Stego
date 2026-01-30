export namespace log {
	
	export class Entry {
	    id: number;
	    // Go type: time
	    timestamp: any;
	    level: string;
	    module: string;
	    message: string;
	    details?: string;
	
	    static createFrom(source: any = {}) {
	        return new Entry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.level = source["level"];
	        this.module = source["module"];
	        this.message = source["message"];
	        this.details = source["details"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace models {
	
	export class AppInfo {
	    name: string;
	    version: string;
	    buildDate: string;
	    buildHash: string;
	    author: string;
	    github: string;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.buildDate = source["buildDate"];
	        this.buildHash = source["buildHash"];
	        this.author = source["author"];
	        this.github = source["github"];
	    }
	}
	export class DecryptRequest {
	    imagePath: string;
	    outputDir: string;
	    password: string;
	    identifier: string;
	
	    static createFrom(source: any = {}) {
	        return new DecryptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.imagePath = source["imagePath"];
	        this.outputDir = source["outputDir"];
	        this.password = source["password"];
	        this.identifier = source["identifier"];
	    }
	}
	export class EncryptRequest {
	    dataSourcePath: string;
	    carrierDir: string;
	    carrierImagePath: string;
	    outputDir: string;
	    outputFileName: string;
	    password: string;
	    scatter?: boolean;
	    identifier: string;
	    autoSelectCarrier: boolean;
	    preferLargestImage: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EncryptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dataSourcePath = source["dataSourcePath"];
	        this.carrierDir = source["carrierDir"];
	        this.carrierImagePath = source["carrierImagePath"];
	        this.outputDir = source["outputDir"];
	        this.outputFileName = source["outputFileName"];
	        this.password = source["password"];
	        this.scatter = source["scatter"];
	        this.identifier = source["identifier"];
	        this.autoSelectCarrier = source["autoSelectCarrier"];
	        this.preferLargestImage = source["preferLargestImage"];
	    }
	}
	export class GenerateRequest {
	    outputDir: string;
	    targetBytes: number;
	    count: number;
	    prefix: string;
	    randomSeed: number;
	    noiseEnabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GenerateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outputDir = source["outputDir"];
	        this.targetBytes = source["targetBytes"];
	        this.count = source["count"];
	        this.prefix = source["prefix"];
	        this.randomSeed = source["randomSeed"];
	        this.noiseEnabled = source["noiseEnabled"];
	    }
	}

}

