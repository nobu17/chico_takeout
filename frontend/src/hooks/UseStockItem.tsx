import { useEffect, useState } from "react";
import { ItemKind } from "../libs/ItemKind";
import { StockItem } from "../libs/StockItem";
import StockItemApi, { StockItemRemainUpdateRequest } from "../libs/apis/stockItem";
import ItemKindApi from "../libs/apis/itemKind";

const defaultStockItem: StockItem = {
  id: "",
  name: "",
  description: "",
  priority: 1,
  price: 100,
  maxOrder: 5,
  enabled: true,
  remain: 0,
  kind: null
};

const stockApi = new StockItemApi("http://localhost:8086");
const kindApi = new ItemKindApi("http://localhost:8086");

export default function useStockItem() {
  const [stockItems, setStockItems] = useState<StockItem[]>([]);
  const [itemKinds, setItemKinds] = useState<ItemKind[]>([]);
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  const addStockItem = (item: StockItem) => {
    add(item);
  };

  const updateStockItem = (item: StockItem) => {
    update(item);
  };

  const deleteStockItem = (item: StockItem) => {
    deletes(item);
  };

  const add = async (item: StockItem) => {
    try {
        setError(undefined);
        setLoading(true);
        await stockApi.add(item);
        // reload
        await getForReload();
      } catch (e: any) {
        setError(e);
      } finally {
        setLoading(false);
      }
  };

  const update = async (item: StockItem) => {
    try {
        setError(undefined);
        setLoading(true);
        await stockApi.update(item);
        // reload
        await getForReload();
      } catch (e: any) {
        setError(e);
      } finally {
        setLoading(false);
      }
  };

  const updateRemain = async (request: StockItemRemainUpdateRequest) => {
    try {
        setError(undefined);
        setLoading(true);
        await stockApi.updateRemain(request);
        // reload
        await getForReload();
      } catch (e: any) {
        setError(e);
      } finally {
        setLoading(false);
      }
  };

  const deletes = async (item: StockItem) => {
    try {
        setError(undefined);
        setLoading(true);
        await stockApi.delete(item);
        // reload
        await getForReload();
      } catch (e: any) {
        setError(e);
      } finally {
        setLoading(false);
      }
  };


  const getForReload = async () => {
    const result = await stockApi.getAll();
    setStockItems(result.data);
  };

  const getAll = async () => {
    try {
      setError(undefined);
      setLoading(true);
      const result = await stockApi.getAll();
      setStockItems(result.data);
      const kinds = await kindApi.getAll();
      setItemKinds(kinds.data);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    stockItems,
    itemKinds,
    defaultStockItem,
    addStockItem,
    updateStockItem,
    deleteStockItem,
    updateRemain,
    error,
    loading,
  };
}
