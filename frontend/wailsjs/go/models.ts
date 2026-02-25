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
	
	export class StorageMigrationResult {
	    from_dir: string;
	    to_dir: string;
	    backed_up_conflicts: string[];
	    restart_recommended: boolean;
	
	    static createFrom(source: any = {}) {
	        return new StorageMigrationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from_dir = source["from_dir"];
	        this.to_dir = source["to_dir"];
	        this.backed_up_conflicts = source["backed_up_conflicts"];
	        this.restart_recommended = source["restart_recommended"];
	    }
	}
	export class StorageSettings {
	    current_data_dir: string;
	    default_data_dir: string;
	    is_custom: boolean;
	    startup_notice: string;
	
	    static createFrom(source: any = {}) {
	        return new StorageSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current_data_dir = source["current_data_dir"];
	        this.default_data_dir = source["default_data_dir"];
	        this.is_custom = source["is_custom"];
	        this.startup_notice = source["startup_notice"];
	    }
	}
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
	    asset_name: string;
	    asset_size: number;
	    can_auto_update: boolean;
	    auto_update_reason: string;
	
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
	        this.asset_name = source["asset_name"];
	        this.asset_size = source["asset_size"];
	        this.can_auto_update = source["can_auto_update"];
	        this.auto_update_reason = source["auto_update_reason"];
	    }
	}

}

