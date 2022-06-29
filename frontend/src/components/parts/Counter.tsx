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
    if (lessThanMax()) {
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

  const lessThanMax = (): boolean => {
    if (!props.max) {
      return true;
    }
    if (count + 1 < props.max) {
      return true;
    }
    return false;
  };

  return (
    <>
      <Stack direction="row" spacing={2}>
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
      </Stack>
    </>
  );
}
