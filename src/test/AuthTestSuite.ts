import { AuthDb } from "../core/classes/AuthDb";
import Test from "./classes/Test";
import TestSuite from "./classes/TestSuite";
import Assertion from "./lib/Assertion";
import Logger from "./lib/Logger";

export class AuthTestSuite extends TestSuite{
    static tests = [
        new Test(async function authDatabaseConnectTest(){

            const db = new AuthDb();
            await db.open();

            db.close();
        }),

        new Test(async function authDatabaseCreateTablesTest(){
                
                const db = new AuthDb();
                await db.open();
                await db.createTables();

                const users = await db.get().all("SELECT * FROM auth");
                Logger.log("Available users", users.map((u : any) => u.username));

                Assertion.assertSize(3, users, "There should be 2 users in the database.");
                db.close();
        })
    ]

    constructor(){
        super(AuthTestSuite.tests);
    }
}
