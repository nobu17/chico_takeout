import { useState } from "react";
import { useAuth } from "../components/contexts/AuthContext";
import StatisticsApi from "../libs/apis/statistics";
import { MonthlyStatisticData } from "../libs/Statistics";
import { getNowDateByMonthOffset } from "../libs/util/DateUtil";

const api = new StatisticsApi();

const defaultData = {
  data: [],
};

export function useMonthlyStatistics() {
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(false);
  const { state } = useAuth();
  const [statistics, setStatistics] =
    useState<MonthlyStatisticData>(defaultData);

  const load = async () => {
    if (!state.isAdmin) {
      setError(new Error("管理ユーザで認証されていません。"));
      return;
    }
    try {
      setError(undefined);
      setLoading(true);
      // 1 years from now
      const req = {
        start: getNowDateByMonthOffset(-12),
        end: getNowDateByMonthOffset(0),
      };
      const result = await api.getMonthly(req);
      const res = result.data;
      if (res && res.data.length > 0) {
        setStatistics(res);
      } else {
        setStatistics(defaultData);
      }
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    statistics,
    load,
    error,
    loading,
  };
}
