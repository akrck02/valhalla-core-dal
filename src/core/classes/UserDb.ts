import Paths from "../config/Paths";
import { Database } from "./Database";

export class UserDb extends Database {

    constructor(name : string) {
        super(name + "-user");
    }

    async createTables() {
        const fs = require('fs');
        const query = fs.readFileSync(Paths.USER_CREATE_TABLES, 'utf8');
        await this.db.exec(query);
    }
}