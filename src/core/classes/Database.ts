import sqlite3 from 'sqlite3'
import { open } from 'sqlite'
import Paths from '../config/Paths';


export class Database {
    protected db : any;
    protected name : string;

    public constructor(name : string) {    
        this.name = name;
    }

    /**
     * Create database
     */
    async open(){
        this.db = await open({
            filename: `${Paths.DATABASE_STORAGE}${this.name}${Paths.DATABASE_EXTENSION}`,
            driver: sqlite3.Database
        });
      
        this.db.on('trace', (data : string) => {
            if(data.indexOf("error") > -1)
                this.log("[SQLITE]", data);
        })

        this.log(`Connected to the in-memory SQLite database ${this.name}.`);
    
    };

    /**
     * Close the database connection
     */
    close() {
        this.db.close((err: Error) => {
            if (err) {
                console.error(err.message);
                return;
            }
            this.log("Closed the database connection.");
        });
    }

    /**
     * Log a message as database
     * @param mgs 
     */
    log(...mgs: any) {
        console.log("Database", mgs);
    }


    get() : any{
        return this.db;
    }

}