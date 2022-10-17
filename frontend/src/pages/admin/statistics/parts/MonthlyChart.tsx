import * as React from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Brush,
} from "recharts";
import { Paper, Select, MenuItem, Box, SelectChangeEvent } from "@mui/material";
import { MonthlyStatisticData } from "../../../../libs/Statistics";

type DisplayType = "Order" | "Quantity" | "Money";

type MonthlyChartProps = {
  data?: MonthlyStatisticData;
};

type RowData = {
  name: string;
  value: number;
};

const createData = (
  data: MonthlyStatisticData,
  displayType: DisplayType
): RowData[] => {
  const rows: Array<RowData> = [];

  for (const item of data.data) {
    switch (displayType) {
      case "Order":
        rows.push({ name: item.month, value: item.orderTotal });
        break;
      case "Quantity":
        rows.push({ name: item.month, value: item.quantityTotal });
        break;
      case "Money":
        rows.push({ name: item.month, value: item.moneyTotal });
        break;
    }
  }

  return rows;
};

export default function MonthlyChart(props: MonthlyChartProps) {
  const [mode, setMode] = React.useState<DisplayType>("Order");
  const handleModeChange = (event: SelectChangeEvent) => {
    setMode(event.target.value as DisplayType);
  };
  if (!props.data) {
    return <></>;
  }
  const data = createData(props.data, mode);
  return (
    <Paper sx={{ m: { md: 1 }, p: { md: 1 } }}>
      <Select
        sx={{ my: 1 }}
        fullWidth
        value={mode}
        label="表示項目"
        onChange={handleModeChange}
      >
        <MenuItem value={"Order"}>注文数</MenuItem>
        <MenuItem value={"Quantity"}>注文商品数</MenuItem>
        <MenuItem value={"Money"}>金額</MenuItem>
      </Select>
      <Box sx={{ height: { xs: 300, md: 600 } }} width="99%">
        <ResponsiveContainer>
          <LineChart
            data={data}
            margin={{
              top: 5,
              right: 30,
              left: 10,
              bottom: 5,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" />
            <YAxis yAxisId="left" />
            <Tooltip />
            <Legend />
            <Brush
              dataKey="name"
              stroke="#8884d8"
              height={30}
              startIndex={0}
              endIndex={data.length - 1}
            />
            <Line
              yAxisId="left"
              type="monotone"
              dataKey="value"
              stroke="#8884d8"
              activeDot={{ r: 8 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </Box>
    </Paper>
  );
}
