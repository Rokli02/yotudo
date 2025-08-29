export namespace database {
	
	export class Database {
	    // Go type: sql
	    Conn?: any;
	
	    static createFrom(source: any = {}) {
	        return new Database(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Conn = this.convertValues(source["Conn"], null);
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

export namespace model {
	
	export class Author {
	    Id: number;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new Author(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	    }
	}
	export class Genre {
	    Id: number;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new Genre(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	    }
	}
	export class Music {
	    Id: number;
	    Name: string;
	    Published?: number;
	    Album?: string;
	    Url: string;
	    Filename?: string;
	    PicFilename?: string;
	    Status: number;
	    Genre: Genre;
	    Author: Author;
	    Contributors: Author[];
	
	    static createFrom(source: any = {}) {
	        return new Music(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Published = source["Published"];
	        this.Album = source["Album"];
	        this.Url = source["Url"];
	        this.Filename = source["Filename"];
	        this.PicFilename = source["PicFilename"];
	        this.Status = source["Status"];
	        this.Genre = this.convertValues(source["Genre"], Genre);
	        this.Author = this.convertValues(source["Author"], Author);
	        this.Contributors = this.convertValues(source["Contributors"], Author);
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
	export class OptionalAuthor {
	    Id?: number;
	    Name?: string;
	
	    static createFrom(source: any = {}) {
	        return new OptionalAuthor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	    }
	}
	export class NewMusic {
	    Name: string;
	    Published: number;
	    Album: string;
	    Url: string;
	    Author: OptionalAuthor;
	    Contributors: OptionalAuthor[];
	    GenreId: number;
	    PicFilename: string;
	
	    static createFrom(source: any = {}) {
	        return new NewMusic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Published = source["Published"];
	        this.Album = source["Album"];
	        this.Url = source["Url"];
	        this.Author = this.convertValues(source["Author"], OptionalAuthor);
	        this.Contributors = this.convertValues(source["Contributors"], OptionalAuthor);
	        this.GenreId = source["GenreId"];
	        this.PicFilename = source["PicFilename"];
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
	
	export class Page {
	    Page: number;
	    Size: number;
	
	    static createFrom(source: any = {}) {
	        return new Page(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Page = source["Page"];
	        this.Size = source["Size"];
	    }
	}
	export class Pagination___yotudo_src_model_Author_ {
	    Data: Author[];
	    Count: number;
	
	    static createFrom(source: any = {}) {
	        return new Pagination___yotudo_src_model_Author_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], Author);
	        this.Count = source["Count"];
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
	export class Pagination___yotudo_src_model_Music_ {
	    Data: Music[];
	    Count: number;
	
	    static createFrom(source: any = {}) {
	        return new Pagination___yotudo_src_model_Music_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], Music);
	        this.Count = source["Count"];
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
	export class Sort {
	    Key: string;
	    Dir: number;
	
	    static createFrom(source: any = {}) {
	        return new Sort(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Key = source["Key"];
	        this.Dir = source["Dir"];
	    }
	}
	export class Status {
	    Id: number;
	    Name: string;
	    Description: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	    }
	}
	export class UpdateMusic {
	    Id: number;
	    Name: string;
	    Published: number;
	    Album: string;
	    Url: string;
	    Author: OptionalAuthor;
	    Contributors: OptionalAuthor[];
	    Status: number;
	    GenreId: number;
	    Filename: string;
	    PicFilename: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateMusic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Published = source["Published"];
	        this.Album = source["Album"];
	        this.Url = source["Url"];
	        this.Author = this.convertValues(source["Author"], OptionalAuthor);
	        this.Contributors = this.convertValues(source["Contributors"], OptionalAuthor);
	        this.Status = source["Status"];
	        this.GenreId = source["GenreId"];
	        this.Filename = source["Filename"];
	        this.PicFilename = source["PicFilename"];
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

export namespace service {
	
	export class InfoService {
	
	
	    static createFrom(source: any = {}) {
	        return new InfoService(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

