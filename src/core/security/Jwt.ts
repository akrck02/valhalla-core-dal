import jwt from 'jsonwebtoken';

export interface JWT {
    username: string;
    password: string;
    mail: string;
}

export function signJWT(data: any, key: string) : string {
    return jwt.sign(data, key);
}

export function verifyJWT(token: string, key: string) : JWT {
    return jwt.verify(token, key) as JWT; 
}
