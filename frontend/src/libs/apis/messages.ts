import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { Message } from "../Messages";

const apiRoot = "/message/store/";

export default class MessagesApi extends ApiBase {
  async get(id: string): Promise<ApiResponse<Message>> {
    const result = await this.getAsync<Message>(apiRoot + id);
    return result;
  }

  async update(item: Message): Promise<void> {
    await this.putAsync(apiRoot + item.id, item);
  }
}
