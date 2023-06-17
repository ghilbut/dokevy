/* eslint-disable */
import {
  CallOptions,
  ChannelCredentials,
  Client,
  ClientOptions,
  ClientUnaryCall,
  handleUnaryCall,
  makeGenericClientConstructor,
  Metadata,
  ServiceError,
  UntypedServiceImplementation,
} from "@grpc/grpc-js";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "dokevy.v1";

export interface Pong {
  pong: string;
}

function createBasePong(): Pong {
  return { pong: "" };
}

export const Pong = {
  encode(message: Pong, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pong !== "") {
      writer.uint32(10).string(message.pong);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Pong {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePong();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.pong = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Pong {
    return { pong: isSet(object.pong) ? String(object.pong) : "" };
  },

  toJSON(message: Pong): unknown {
    const obj: any = {};
    message.pong !== undefined && (obj.pong = message.pong);
    return obj;
  },

  create<I extends Exact<DeepPartial<Pong>, I>>(base?: I): Pong {
    return Pong.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Pong>, I>>(object: I): Pong {
    const message = createBasePong();
    message.pong = object.pong ?? "";
    return message;
  },
};

export type SystemServiceService = typeof SystemServiceService;
export const SystemServiceService = {
  ping: {
    path: "/dokevy.v1.SystemService/Ping",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: Empty) => Buffer.from(Empty.encode(value).finish()),
    requestDeserialize: (value: Buffer) => Empty.decode(value),
    responseSerialize: (value: Pong) => Buffer.from(Pong.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Pong.decode(value),
  },
} as const;

export interface SystemServiceServer extends UntypedServiceImplementation {
  ping: handleUnaryCall<Empty, Pong>;
}

export interface SystemServiceClient extends Client {
  ping(request: Empty, callback: (error: ServiceError | null, response: Pong) => void): ClientUnaryCall;
  ping(
    request: Empty,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Pong) => void,
  ): ClientUnaryCall;
  ping(
    request: Empty,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Pong) => void,
  ): ClientUnaryCall;
}

export const SystemServiceClient = makeGenericClientConstructor(
  SystemServiceService,
  "dokevy.v1.SystemService",
) as unknown as {
  new (address: string, credentials: ChannelCredentials, options?: Partial<ClientOptions>): SystemServiceClient;
  service: typeof SystemServiceService;
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
