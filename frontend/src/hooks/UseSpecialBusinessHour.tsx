import { useEffect, useState } from "react";
import { SpecialBusinessHour } from "../libs/SpecialBusinessHour";
import SpecialBusinessHourApi from "../libs/apis/specialBusinessHour";
import { toDateString, getNowDate } from "../libs/util/DateUtil";
import { ApiError } from "../libs/apis/apibase";

const defaultSpecialBusinessHour: SpecialBusinessHour = {
  id: "",
  name: "",
  start: "08:00",
  end: "10:00",
  date: toDateString(getNowDate(0)),
  businessHourId: "",
  offsetHour: 3,
};

const busApi = new SpecialBusinessHourApi();

export default function useSpecialBusinessHour() {
  const [specialBusinessHours, setSpecialBusinessHours] = useState<SpecialBusinessHour[]>([]);
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  const addSpecialBusinessHour = (item: SpecialBusinessHour) : Promise<ApiError | null> => {
    return add(item);
  };

  const add = async (item: SpecialBusinessHour): Promise<ApiError | null> => {
    try {
        setError(undefined);
        setLoading(true);
        await busApi.add(item);
        // reload
        await getForReload();
        return null;
      } catch (e: any) {
        // bad request not set error
        if (e instanceof ApiError) {
          if (!e.isBadRequest()){
            setError(e);
          }
        } 
        return e;
      } finally {
        setLoading(false);
      }
  };


  const updateSpecialBusinessHour = (item: SpecialBusinessHour) : Promise<ApiError | null> => {
    return update(item);
  };

  const update = async (item: SpecialBusinessHour): Promise<ApiError | null> => {
    try {
        setError(undefined);
        setLoading(true);
        await busApi.update(item);
        // reload
        await getForReload();
        return null;
      } catch (e: any) {
        // bad request not set error
        if (e instanceof ApiError) {
          if (!e.isBadRequest()){
            setError(e);
          }
        } 
        return e;
      } finally {
        setLoading(false);
      }
  };

  const getForReload = async () => {
    const result = await busApi.getAll();
    setSpecialBusinessHours(result.data);
  };

  const deleteSpecialBusinessHour = (item: SpecialBusinessHour) => {
      deletes(item);
  };

  const deletes = async (item: SpecialBusinessHour) => {
    try {
      setError(undefined);
      setLoading(true);
      await busApi.delete(item);
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
      const result = await busApi.getAll();
      setSpecialBusinessHours(result.data);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    specialBusinessHours,
    addSpecialBusinessHour,
    updateSpecialBusinessHour,
    deleteSpecialBusinessHour,
    defaultSpecialBusinessHour,
    error,
    loading,
  };
}
