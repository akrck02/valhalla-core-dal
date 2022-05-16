import Paths from "../config/Paths";
import { Database } from "./Database";

export class AuthDb extends Database {

    constructor() {
        super("auth");
    }

    async createTables() {
        const fs = require('fs');
        const query = fs.readFileSync(Paths.AUTH_CREATE_TABLES, 'utf8');
        await this.db.exec(query);
    } 
}