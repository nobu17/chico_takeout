import { useEffect, useState } from "react";
import { ItemKind } from "../libs/ItemKind";
import ItemKindApi from "../libs/apis/itemKind";

const defaultItemKind: ItemKind = {
  id: "",
  name: "",
  priority: 1,
};

const api = new ItemKindApi();

export default function useItemKind() {
  const [itemKinds, setItemKinds] = useState<ItemKind[]>([]);
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  const addNewItemKind = (item: ItemKind) => {
    add(item);
  };

  const updateItemKind = (item: ItemKind) => {
    update(item);
  };

  const deleteItemKind = (item: ItemKind) => {
    deletes(item);
  };

  const add = async (item: ItemKind) => {
    try {
      setError(undefined);
      setLoading(true);
      await api.add(item);
      // reload
      await getForReload();
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const update = async (item: ItemKind) => {
    try {
      setError(undefined);
      setLoading(true);
      await api.update(item);
      // reload
      await getForReload();
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const deletes = async (item: ItemKind) => {
    try {
      setError(undefined);
      setLoading(true);
      await api.delete(item);
      // reload
      await getForReload();
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const getForReload = async () => {
    const result = await api.getAll();
    setItemKinds(result.data);
  };

  const getAll = async () => {
    try {
      setError(undefined);
      setLoading(true);
      const result = await api.getAll();
      setItemKinds(result.data);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    itemKinds,
    defaultItemKind,
    addNewItemKind,
    updateItemKind,
    deleteItemKind,
    error,
    loading,
  };
}
