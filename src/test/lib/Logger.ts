export default class Logger {

    static log(...data : any) {
        console.info("   ",data.join(","));
    }

    static success(...data : any) {
        console.info(" ✅", data.join(","));
    }

    static warning(...data : any) {
        console.warn(" 🔶", data.join(","));
    }

    static error(...data : any) {
        console.error(" ❌", data.join(","));
    }

    static jump() {
        console.info("");
    }

    static rawlog(...data : any) {
        console.info(data.join(" "));
    }

    static hardTitle(title : string) {
        console.info("\n##########################################");
        console.info("  " + title.toUpperCase() + " ");
        console.info("##########################################");
    }

    static softTitle(title : string) {
        console.info("\n------------------------------------------");
        console.info("  " + title);
        console.info("------------------------------------------");
    }
}