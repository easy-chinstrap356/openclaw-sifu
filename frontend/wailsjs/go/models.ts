export namespace main {
	
	export class Environment {
	    hostname: string;
	    platform: string;
	    architecture: string;
	
	    static createFrom(source: any = {}) {
	        return new Environment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.platform = source["platform"];
	        this.architecture = source["architecture"];
	    }
	}
	export class BootstrapPayload {
	    environment: Environment;
	
	    static createFrom(source: any = {}) {
	        return new BootstrapPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.environment = this.convertValues(source["environment"], Environment);
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
	
	export class InstallerConfig {
	    tag: string;
	    installMethod: string;
	    gitDir: string;
	    noOnboard: boolean;
	    noGitUpdate: boolean;
	    dryRun: boolean;
	    useCnMirrors: boolean;
	    npmRegistry: string;
	    installBaseUrl: string;
	    repoUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new InstallerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag = source["tag"];
	        this.installMethod = source["installMethod"];
	        this.gitDir = source["gitDir"];
	        this.noOnboard = source["noOnboard"];
	        this.noGitUpdate = source["noGitUpdate"];
	        this.dryRun = source["dryRun"];
	        this.useCnMirrors = source["useCnMirrors"];
	        this.npmRegistry = source["npmRegistry"];
	        this.installBaseUrl = source["installBaseUrl"];
	        this.repoUrl = source["repoUrl"];
	    }
	}
	export class InstallerResult {
	    success: boolean;
	    installedVersion: string;
	    isUpgrade: boolean;
	    message: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new InstallerResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.installedVersion = source["installedVersion"];
	        this.isUpgrade = source["isUpgrade"];
	        this.message = source["message"];
	        this.error = source["error"];
	    }
	}
	export class PostInstallActionResult {
	    success: boolean;
	    message: string;
	    error?: string;
	    cancelled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PostInstallActionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.error = source["error"];
	        this.cancelled = source["cancelled"];
	    }
	}

}

