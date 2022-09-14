import { useEffect, useState } from "react";
import { BusinessHour } from "../libs/BusinessHour";
import businessHourApi from "../libs/apis/businessHour";
import { ApiError } from "../libs/apis/apibase";

const defaultBusinessHour: BusinessHour = {
  id: "",
  name: "",
  start: "08:00",
  end: "10:00",
  weekdays: [],
};

const busApi = new businessHourApi();

export default function useFoodItem() {
  const [businessHours, setBusinessHours] = useState<BusinessHour[]>([]);
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  const updateBusinessHour = (item: BusinessHour) : Promise<ApiError | null> => {
    return update(item);
  };

  const update = async (item: BusinessHour): Promise<ApiError | null> => {
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
    setBusinessHours(result.data.schedules);
  };

  const getAll = async () => {
    try {
      setError(undefined);
      setLoading(true);
      const result = await busApi.getAll();
      setBusinessHours(result.data.schedules);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    businessHours,
    defaultBusinessHour,
    updateBusinessHour,
    error,
    loading,
  };
}
