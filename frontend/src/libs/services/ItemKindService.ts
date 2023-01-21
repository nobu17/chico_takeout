import ItemKindApi from "../apis/itemKind";
import OptionItemApi from "../apis/optionItem";
import { AggregatedItemKind } from "../ItemKind";

export default class ItemKindService {
  private kindApi: ItemKindApi;
  private optionItemApi: OptionItemApi;

  constructor(public baseUrl: string = "") {
    this.kindApi = new ItemKindApi(baseUrl);
    this.optionItemApi = new OptionItemApi(baseUrl);
  }

  async getAggregatedItemKinds(): Promise<AggregatedItemKind[]> {
    const kinds = await this.kindApi.getAll();
    const optionItems = await this.optionItemApi.getAll();

    const aggregatedKinds: AggregatedItemKind[] = [];
    for (const kind of kinds.data) {
      let agKind = {
        id: kind.id,
        name: kind.name,
        priority: kind.priority,
        options: optionItems.data.filter(
          (x) => x.enabled && kind.optionItemIds.includes(x.id)
        ),
      };
      aggregatedKinds.push(agKind);
    }

    return aggregatedKinds;
  }
}
