export class DataCache<T = unknown> {
    private ttl: number;
    private lastModified: number;
    private intervalId: any;
    private _data?: T;

    /**
     * @param ttl Object's Time To Live measured in seconds.
     * @param _data Data which gets stored in this cache.
     */
    constructor(_data?: T, ttl?: number) {
        // 30 seconds
        this.ttl = ttl ? ttl * 1000 : 30000
        this.lastModified = Date.now();

        if (_data != undefined) {
            this.data = _data
        }
    }

    
    public get data() : T | undefined {
        return this._data;
    }
    
    
    public set data(value : T | undefined) {
        if (this.intervalId !== undefined) {
            this.stopInterval()
        }
        this._data = value;
        
        if (value !== undefined) {
            this.lastModified = Date.now();
            this.startInterval();
        }
    }

    private startInterval() {
        const context = this

        this.intervalId = setInterval(() => {
            if (context.lastModified + context.ttl < Date.now()) {
                context._data = undefined;
                context.stopInterval()
            }
        }, this.ttl)
    }

    private stopInterval() {
        clearInterval(this.intervalId)
        this.intervalId = undefined
    }
}