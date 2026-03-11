export namespace main {
	
	export class AgentProfile {
	    name: string;
	    summary: string;
	    endpointMode: string;
	    promptLayer: string;
	
	    static createFrom(source: any = {}) {
	        return new AgentProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.summary = source["summary"];
	        this.endpointMode = source["endpointMode"];
	        this.promptLayer = source["promptLayer"];
	    }
	}
	export class PipelineStep {
	    title: string;
	    description: string;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new PipelineStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.description = source["description"];
	        this.state = source["state"];
	    }
	}
	export class Capability {
	    id: string;
	    title: string;
	    detail: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Capability(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.detail = source["detail"];
	        this.status = source["status"];
	    }
	}
	export class Environment {
	    hostname: string;
	    username: string;
	    platform: string;
	    architecture: string;
	    goVersion: string;
	    workingDir: string;
	    executablePath: string;
	    tempDir: string;
	    powerShellPath: string;
	    webView2State: string;
	
	    static createFrom(source: any = {}) {
	        return new Environment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.username = source["username"];
	        this.platform = source["platform"];
	        this.architecture = source["architecture"];
	        this.goVersion = source["goVersion"];
	        this.workingDir = source["workingDir"];
	        this.executablePath = source["executablePath"];
	        this.tempDir = source["tempDir"];
	        this.powerShellPath = source["powerShellPath"];
	        this.webView2State = source["webView2State"];
	    }
	}
	export class BootstrapPayload {
	    appName: string;
	    version: string;
	    modeLabel: string;
	    summary: string;
	    bootTime: string;
	    agent: AgentProfile;
	    environment: Environment;
	    capabilities: Capability[];
	    pipeline: PipelineStep[];
	    nextSteps: string[];
	
	    static createFrom(source: any = {}) {
	        return new BootstrapPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appName = source["appName"];
	        this.version = source["version"];
	        this.modeLabel = source["modeLabel"];
	        this.summary = source["summary"];
	        this.bootTime = source["bootTime"];
	        this.agent = this.convertValues(source["agent"], AgentProfile);
	        this.environment = this.convertValues(source["environment"], Environment);
	        this.capabilities = this.convertValues(source["capabilities"], Capability);
	        this.pipeline = this.convertValues(source["pipeline"], PipelineStep);
	        this.nextSteps = source["nextSteps"];
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

}

