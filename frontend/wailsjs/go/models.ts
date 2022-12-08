export namespace main {
	
	export class Message {
	    topic: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.topic = source["topic"];
	        this.message = source["message"];
	    }
	}

}

export namespace model {
	
	export class ChangePasswordParam {
	    username: string;
	    oldPassword: string;
	    newPassword: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangePasswordParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.oldPassword = source["oldPassword"];
	        this.newPassword = source["newPassword"];
	    }
	}

}

