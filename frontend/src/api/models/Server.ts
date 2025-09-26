import { model } from "@wailsjs/go/models";

export interface ServerConfig {
    port: number;
    isHosted: boolean;
    listeningOn: string;
}

export class ServerConfigParser {
    static fromGo(config: model.ServerConfig): ServerConfig {
        return {
            port: config.Port,
            isHosted: config.IsHosted,
            listeningOn: config.ListeningOn,
        };
    }
    static toGo(config: ServerConfig): model.ServerConfig {
        return model.ServerConfig.createFrom({
            Port: config.port,
            IsHosted: config.isHosted,
            ListeningOn: config.listeningOn,
        } as model.ServerConfig)
    }
}