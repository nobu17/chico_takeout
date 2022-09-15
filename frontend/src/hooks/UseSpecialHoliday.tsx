import { useEffect, useState } from "react";
import { SpecialHoliday } from "../libs/SpecialHoliday";
import SpecialHolidayApi from "../libs/apis/specialHoliday";
import { toDateString, getNowDate } from "../libs/util/DateUtil";
import { ApiError } from "../libs/apis/apibase";

const defaultSpecialHoliday: SpecialHoliday = {
  id: "",
  name: "",
  start: toDateString(getNowDate(7)),
  end: toDateString(getNowDate(8)),
};

const api = new SpecialHolidayApi();

export default function useSpecialBusinessHour() {
  const [specialHolidays, setSpecialHolidays] = useState<SpecialHoliday[]>([]);
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  const addSpecialHoliday = (
    item: SpecialHoliday
  ): Promise<ApiError | null> => {
    return add(item);
  };

  const add = async (item: SpecialHoliday): Promise<ApiError | null> => {
    try {
      setError(undefined);
      setLoading(true);
      await api.add(item);
      // reload
      await getForReload();
      return null;
    } catch (e: any) {
      // bad request not set error
      if (e instanceof ApiError) {
        if (!e.isBadRequest()) {
          setError(e);
        }
      }
      return e;
    } finally {
      setLoading(false);
    }
  };

  const updateSpecialHoliday = (
    item: SpecialHoliday
  ): Promise<ApiError | null> => {
    return update(item);
  };

  const update = async (item: SpecialHoliday): Promise<ApiError | null> => {
    try {
      setError(undefined);
      setLoading(true);
      await api.update(item);
      // reload
      await getForReload();
      return null;
    } catch (e: any) {
      // bad request not set error
      if (e instanceof ApiError) {
        if (!e.isBadRequest()) {
          setError(e);
        }
      }
      return e;
    } finally {
      setLoading(false);
    }
  };

  const getForReload = async () => {
    const result = await api.getAll();
    setSpecialHolidays(result.data);
  };

  const deleteSpecialHoliday = (item: SpecialHoliday) => {
    deletes(item);
  };

  const deletes = async (item: SpecialHoliday) => {
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

  const getAll = async () => {
    try {
      setError(undefined);
      setLoading(true);
      const result = await api.getAll();
      setSpecialHolidays(result.data);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    specialHolidays,
    addSpecialHoliday,
    updateSpecialHoliday,
    deleteSpecialHoliday,
    defaultSpecialHoliday,
    error,
    loading,
  };
}
