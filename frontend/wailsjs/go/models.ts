export namespace database {
	
	export class Database {
	
	
	    static createFrom(source: any = {}) {
	        return new Database(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace main {
	
	export class UpdateCheckResult {
	    has_update: boolean;
	    current_version: string;
	    latest_version: string;
	    release_name: string;
	    release_url: string;
	    published_at: string;
	    release_notes: string;
	    checked_at: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.has_update = source["has_update"];
	        this.current_version = source["current_version"];
	        this.latest_version = source["latest_version"];
	        this.release_name = source["release_name"];
	        this.release_url = source["release_url"];
	        this.published_at = source["published_at"];
	        this.release_notes = source["release_notes"];
	        this.checked_at = source["checked_at"];
	        this.message = source["message"];
	    }
	}

}

