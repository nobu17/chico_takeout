import { useEffect, useState } from "react";
import { OptionItem } from "../libs/OptionItem";
import OptionItemApi from "../libs/apis/optionItem";

const defaultOptionItem: OptionItem = {
  id: "",
  name: "",
  priority: 1,
  description: "",
  price: 100,
  enabled: true,
};

const api = new OptionItemApi();

export default function useOptionItem() {
  const [optionItems, setOptionItems] = useState<OptionItem[]>([]);
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  const addNewOptionItem = (item: OptionItem) => {
    add(item);
  };

  const updateOptionItem = (item: OptionItem) => {
    update(item);
  };

  const deleteOptionItem = (item: OptionItem) => {
    deletes(item);
  };

  const add = async (item: OptionItem) => {
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

  const update = async (item: OptionItem) => {
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

  const deletes = async (item: OptionItem) => {
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
    setOptionItems(result.data);
  };

  const getAll = async () => {
    try {
      setError(undefined);
      setLoading(true);
      const result = await api.getAll();
      setOptionItems(result.data);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    optionItems,
    defaultOptionItem,
    addNewOptionItem,
    updateOptionItem,
    deleteOptionItem,
    error,
    loading,
  };
}
