'use server';

import { NextResponse } from 'next/server'
import {credentials, ServiceError} from "@grpc/grpc-js";
import { Empty } from '~/grpc/v1/google/protobuf/empty';
import {Pong, SystemServiceClient} from "~/grpc/v1/system";

const client = new SystemServiceClient(
    'localhost:50051',
    credentials.createInsecure()
);

const empty: Empty = {
};

const ping = async (): Promise<string> => {

    return new Promise(async (resolve, reject) => {

        const callback = (err: ServiceError | null, response: Pong) => {
            if (err != null) {
                console.error(err);
                reject(err);
            } else {
                console.log("Pong: ", response);
                resolve(response.pong);
            }
        };

        try {
            const call = await client.ping(empty, callback);
            console.log("call: ", call);
        } catch (e) {
            console.error(e);
        }

    });
}

export async function GET() {
    const body = await ping();
    return NextResponse.json(body)
}
