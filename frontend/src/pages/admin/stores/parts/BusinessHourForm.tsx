import * as React from "react";
import { Button, Container, Stack, TextField } from "@mui/material";
import { RhfTimeSelect } from "../../../../components/parts/Rhf/RhfTimeSelect";
import { RhfDayOfWeekSelect } from "../../../../components/parts/Rhf/RhfDayOfWeekSelect";
import { BusinessHour } from "../../../../libs/BusinessHour";
import { StoreTimeList } from "../../../../libs/Constant";
import { SubmitHandler, useForm, FieldError } from "react-hook-form";

import { startIsLessThanEnd } from "../../../../libs/util/TimeCompare";

type BusinessHourFormProps = {
  editItem: BusinessHour;
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (itemKind: BusinessHour): void;
}
interface callbackCancel {
  (): void;
}

const errorMessage = ({
  name,
  error,
  maxLength,
  min,
  max,
}: {
  name: string;
  error: FieldError | undefined;
  maxLength?: string;
  min?: string;
  max?: string;
}): string => {
  if (error?.type === "required") {
    return "入力が必要です。";
  }
  if (error?.type === "maxLength") {
    return maxLength + "文字以下にしてください。";
  }
  if (error?.type === "min") {
    return min + "以上にしてください。";
  }
  if (error?.type === "max") {
    return max + "以下にしてください。";
  }
  return "";
};

export default function BusinessHourForm(props: BusinessHourFormProps) {
  const {
    register,
    control,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<BusinessHour>({ defaultValues: props.editItem });
  const onSubmit: SubmitHandler<BusinessHour> = (data) => {
    if (!startIsLessThanEnd(data.start, data.end)) {
      setError(`start`, { message: "開始時刻が終了時刻よりも大きいです。" });
      setError(`end`, { message: "開始時刻が終了時刻よりも大きいです。" });
      return;
    }
    props.onSubmit(data);
  };
  const onCancel = () => {
    props.onCancel();
  };
  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <TextField
            label="名称"
            {...register("name", { required: true, maxLength: 10 })}
            error={Boolean(errors.name)}
            helperText={errorMessage({
              name: "名称",
              error: errors.name,
              maxLength: "10",
            })}
          />
          <RhfTimeSelect
            label="開始時間"
            timeList={timeList}
            name="start"
            control={control}
          />
          <RhfTimeSelect
            label="終了時間"
            timeList={timeList}
            name="end"
            control={control}
          />
          <RhfDayOfWeekSelect
            label="曜日"
            name="weekdays"
            control={control}
          />
          <Button
            color="primary"
            variant="contained"
            size="large"
            onClick={handleSubmit(onSubmit)}
          >
            確定
          </Button>
          <Button
            color="error"
            variant="contained"
            size="large"
            onClick={onCancel}
          >
            キャンセル
          </Button>
        </Stack>
      </Container>
    </>
  );
}

const timeList = StoreTimeList;