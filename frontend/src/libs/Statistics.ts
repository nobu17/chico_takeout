export type MonthlyStatisticData = {
    data: MonthlyStatisticRecord[],
}

export type MonthlyStatisticRecord = {
  month: string;
  orderTotal: number;
  quantityTotal: number;
  moneyTotal: number;
};
