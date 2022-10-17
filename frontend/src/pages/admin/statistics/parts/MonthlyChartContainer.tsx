import { useEffect } from "react";
import { Alert } from "@mui/material";
import { useMonthlyStatistics } from "../../../../hooks/UseMonthlyStatistics";
import MonthlyChart from "./MonthlyChart";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";

const errorMessage = (error: Error | undefined) => {
  if (error) {
    console.log("err", error);
    return (
      <Alert variant="filled" severity="error">
        エラーが発生しました。
      </Alert>
    );
  }
  return <></>;
};

export default function MonthlyChartContainer() {
  const { statistics, load, error, loading } = useMonthlyStatistics();

  useEffect(() => {
    const init = async () => {
      await load();
    };
    init();
  }, []);

  return (
    <>
      {errorMessage(error)}
      <MonthlyChart data={statistics}></MonthlyChart>
      <LoadingSpinner message="Loading..." isLoading={loading} />
    </>
  );
}
