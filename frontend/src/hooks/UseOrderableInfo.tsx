import { useState, useEffect } from "react";
import OrderableInfoService from "../libs/services/OrderableInfoService";
import { PerDayOrderableInfo } from "../libs/OrderableInfo";
import { IsTimeInRange } from "../libs/util/DateUtil";

const defaultPerDayOrderableInfo: PerDayOrderableInfo[] = [
  {
    date: "",
    hourTypeId: "",
    startTime: "",
    endTime: "",
    categories: [],
  },
];

const service = new OrderableInfoService();

export function useOrderableInfo() {
  const [perDayOrderableInfo, setPerDayOrderableInfo] = useState<
    PerDayOrderableInfo[]
  >(defaultPerDayOrderableInfo);
  const [currentOrderableInfo, setCurrentOrderableInfo] =
    useState<PerDayOrderableInfo>(defaultPerDayOrderableInfo[0]);
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  const getAll = async () => {
    try {
      setError(undefined);
      setLoading(true);
      const result = await service.get();
      setPerDayOrderableInfo(result);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const switchCurrent = (date: string, time: string) => {
    const matched = perDayOrderableInfo.find(
      (p) => p.date === date && IsTimeInRange(p.startTime, p.endTime, time)
    );
    if (matched) {
      setCurrentOrderableInfo(matched);
    }
  };

  return {
    perDayOrderableInfo,
    currentOrderableInfo,
    switchCurrent,
    error,
    loading,
  };
}
