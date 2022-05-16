import { Database } from "sqlite";

export class TaskData {

    /**
     * 
     * @param db The databas3e connection
     * @param username The user who owns the tasks
     * @returns The query result
     */
     public static getUserTasks(db: Database, username: string): Promise<any> {
        const SQL = "SELECT * FROM task WHERE author = ? ORDER BY end DESC";
        const response = db.all(SQL,username);
        return response;
    }
    
}