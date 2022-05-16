import { Database } from "sqlite3";

export class NoteData {

    /**
     * Get the notes of a user
     * @param db The database connection
     * @param username The username to search for
     * @returns The notes of a given user
     */
     public static async getuserNotes(db: Database, username : string ) {
        const SQL = "SELECT * FROM note WHERE author = ?";
        const response = db.all(SQL, username);
        return response;
    }



}