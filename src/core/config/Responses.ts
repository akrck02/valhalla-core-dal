interface HttpResponse {
    success: boolean
    message: string 
    code : number
}

export const PONG = new Promise((r) => r({
    success: true , 
    message : "pong", 
    code: 200}
));

export const MISSING_PARAMETERS = new Promise<HttpResponse>((r) => r({
    success: false, 
    message:"Missing parameters.", 
    code : 400
}));

export const INCORRECT_CREDENTIALS = new Promise<HttpResponse>((r) => r({
    success: false, 
    message:"Incorrect credentials.", 
    code : 403
}));

export const NOT_IMPLEMENTED_YET = new Promise<HttpResponse>((r) => r({
    success: false, 
    message:"Not implemented yet.", 
    code : 404
}));

export const SOMETHING_WENT_WRONG = new Promise<HttpResponse>((r) => r({
    success: false, 
    message: "Something went wrong.", 
    code : 500
})); 


export const MAIL_ALREADY_EXISTS = new Promise<HttpResponse>((r) => r({
    success: false, 
    message: "Mail already exists.", 
    code : 600
})); 

export const MAIL_DOES_NOT_EXIST = new Promise<HttpResponse>((r) => r({
    success: false, 
    message: "Mail does not exist.", 
    code : 601
})); 

export const USER_ALREADY_EXISTS = new Promise<HttpResponse>((r) => r({
    success: false, 
    message: "User already exists.", 
    code : 603
})); 

export const USER_DOES_NOT_EXIST = new Promise<HttpResponse>((r) => r({
    success: false, 
    message: "User does not exist.", 
    code : 604
})); 
