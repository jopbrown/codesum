export namespace cfgs {
	
	export class ChatGpt {
	    end_point?: string;
	    api_key?: string;
	    access_token?: string;
	    model?: string;
	    proxy?: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatGpt(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.end_point = source["end_point"];
	        this.api_key = source["api_key"];
	        this.access_token = source["access_token"];
	        this.model = source["model"];
	        this.proxy = source["proxy"];
	    }
	}
	export class Prompt {
	    system?: string;
	    code_summary?: string;
	    summary_table?: string;
	    final_summary?: string;
	
	    static createFrom(source: any = {}) {
	        return new Prompt(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.system = source["system"];
	        this.code_summary = source["code_summary"];
	        this.summary_table = source["summary_table"];
	        this.final_summary = source["final_summary"];
	    }
	}
	export class SummaryRules {
	    include?: string[];
	    exclude?: string[];
	    out_dir?: string;
	    out_file_name?: string;
	
	    static createFrom(source: any = {}) {
	        return new SummaryRules(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.include = source["include"];
	        this.exclude = source["exclude"];
	        this.out_dir = source["out_dir"];
	        this.out_file_name = source["out_file_name"];
	    }
	}
	export class Config {
	    debug_mode?: boolean;
	    log_path?: string;
	    chat_gpt?: ChatGpt;
	    summary_rules?: SummaryRules;
	    prompt?: Prompt;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.debug_mode = source["debug_mode"];
	        this.log_path = source["log_path"];
	        this.chat_gpt = this.convertValues(source["chat_gpt"], ChatGpt);
	        this.summary_rules = this.convertValues(source["summary_rules"], SummaryRules);
	        this.prompt = this.convertValues(source["prompt"], Prompt);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

export namespace openai {
	
	export class ChatCompletionMessage {
	    role: string;
	    content: string;
	    name?: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatCompletionMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	        this.name = source["name"];
	    }
	}

}

