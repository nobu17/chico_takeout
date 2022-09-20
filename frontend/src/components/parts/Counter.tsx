import * as React from "react";
import { Typography, IconButton, Stack } from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import RemoveIcon from "@mui/icons-material/Remove";

type CounterProps = {
  count: number;
  max?: number;
  onChanged?: (count: number) => void;
};

export default function Counter(props: CounterProps) {
  const [count, setCount] = React.useState(props.count);

  const handleAdd = () => {
    if (lessThanEqualMax()) {
      const newVal = count + 1;
      setCount(newVal);
      props.onChanged?.(newVal);
    }
  };

  const handleRemove = () => {
    if (count > 0) {
      const newVal = count - 1;
      setCount(newVal);
      props.onChanged?.(newVal);
    }
  };

  const lessThanEqualMax = (): boolean => {
    if (!props.max) {
      return true;
    }
    if (count + 1 <= props.max) {
      return true;
    }
    return false;
  };

  const getContent = () => {
    if (props.max !== undefined && props.max <= 0) {
      return (
        <Typography
          align="center"
          color="error"
          gutterBottom
          sx={{ mt: 3, mx: 3 }}
        >
          {"在庫なし"}
        </Typography>
      );
    }
    return (
      <>
        <IconButton color="primary" onClick={handleRemove}>
          <RemoveIcon />
        </IconButton>
        <Typography
          sx={{ py: 1, px: 2, border: 1 }}
          align="justify"
          textAlign="center"
        >
          {count}
        </Typography>
        <IconButton color="error" onClick={handleAdd}>
          <AddIcon />
        </IconButton>
      </>
    );
  };

  return (
    <>
      <Stack direction="row" spacing={2}>
        {getContent()}
        {/* <IconButton color="primary" onClick={handleRemove}>
          <RemoveIcon />
        </IconButton>
        <Typography
          sx={{ py: 1, px: 2, border: 1 }}
          align="justify"
          textAlign="center"
        >
          {count}
        </Typography>
        <IconButton color="error" onClick={handleAdd}>
          <AddIcon />
        </IconButton> */}
      </Stack>
    </>
  );
}
