import React from 'react';

import s from './index.scss';

interface Props {
  isEdit: boolean;
  name: string;
  addText: string;
  updateText: string;
}

const ModalTitle: React.FC<Props> = ({ isEdit, name, addText, updateText }) => {
  return (
    <div className={s.ModalTitleBox}>
      <div className={s.ModalTitleCustom}>{isEdit ? updateText : addText}</div>
      {name}
    </div>
  );
};

export default ModalTitle;
