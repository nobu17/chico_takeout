import React from "react";

type LineBreakProps = {
  msg: string;
};

export default function LineBreak(props: LineBreakProps) {
  const texts = props.msg.split(/(\n)/).map((item, index) => {
    return (
      <React.Fragment key={index}>
        {item.match(/\n/) ? <br /> : item}
      </React.Fragment>
    );
  });
  return <>{texts}</>;
}