import { useState } from "react";
import { useAuth } from "../components/contexts/AuthContext";
import MessagesApi from "../libs/apis/messages";
import { Message } from "../libs/Messages";

const api = new MessagesApi();

export type SelectableData = {
  id: string;
  title: string;
};

const defaultSelectableData = [
  { id: "1", title: "トップページ" },
  { id: "2", title: "マイページ" },
];

export function useMessages() {
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(false);
  const { state } = useAuth();
  const [selectable] = useState<SelectableData[]>(defaultSelectableData);
  const [current, setCurrent] = useState<Message>();

  const load = async (id: string) => {
    if (!selectable.find((s) => s.id === id)) {
      setError(new Error("存在しないIDを選択しました。"));
      return;
    }
    try {
      setError(undefined);
      setLoading(true);
      const result = await api.get(id);
      setCurrent(result.data);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const update = async (updateTarget: Message) => {
    if (!state.isAdmin) {
      setError(new Error("管理ユーザで認証されていません。"));
      return;
    }
    try {
      setError(undefined);
      setLoading(true);
      await api.update(updateTarget);
      // reload
      const result = await api.get(updateTarget.id);
      setCurrent(result.data);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    selectable,
    current,
    load,
    update,
    error,
    loading,
  };
}
