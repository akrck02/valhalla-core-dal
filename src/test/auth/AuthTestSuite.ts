import { AuthDb } from "../../core/classes/AuthDb";
import AuthData from "../../core/data/AuthData";
import Test from "../classes/Test";
import TestSuite from "../classes/TestSuite";
import Assertion from "../lib/Assertion";
import Logger from "../lib/Logger";

export class AuthTestSuite extends TestSuite{
    static tests = [

        /**
         * 
         */
        new Test(async function authDatabaseConnectTest(){

            const db = new AuthDb("auth-test");
            await db.open();

            db.close();
        }),


        /**
         * 
         */
        new Test(async function authDatabaseCreateTablesTest(){
                
            const db = new AuthDb("auth-test-preload");
            await db.open();
            await db.createTables();

            const users = await db.get().all("SELECT * FROM auth");
            Logger.log("Available users", users.map((u : any) => u.username));

            Assertion.assertSize(3, users, "There should be 2 users in the database.");
            db.close();
        }),

        /**
         * 
         */
        new Test(async function authRegisterTest() {
            
            const db = new AuthDb("auth-test");
            await db.open();

            // Create database tables if they don't exist
            await db.createTables();

            // Clean database users
            await db.get().exec("DELETE FROM auth");

            const secret = "SUPER_SECRET_KEY";
            await AuthData.register({
                device : "127.0.0.1",
                user : "dbuser",
                mail : "dbuser@valhalla.com",
                password : "dbuser",
                platform : "Linux"
            },db.get(),secret);

            const users = await db.get().all("SELECT * FROM auth");
            Logger.log("Available users: " + users.map((u : any) => u.username));

            db.close();
        }),

        /**
         * 
         */
        new Test(async function authLoginTest() {
            const db = new AuthDb("auth-test");
            await db.open();

            const login = await AuthData.login({
                device : "127.0.0.1",
                user : "dbuser",
                mail : "dbuser@valhalla.com",
                password : "dbuser",
                platform : "Linux"
            },db.get());            

            Assertion.assert(login, "Login failed");

            db.close();
        })


    ]

    constructor(){
        super(AuthTestSuite.tests);
    }
}
