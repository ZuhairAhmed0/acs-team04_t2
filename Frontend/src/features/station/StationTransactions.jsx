import TransactionsView from '../../components/common/TransactionsView';
import { exportStationTransactions, fetchStationTransactions } from '../../services/station.service';

export default function StationTransactions() {
  return (
    <TransactionsView
      title="كل المعاملات"
      subtitle="كل العمليات المرتبطة بمحطتك، مع إمكانية التصفية وتنزيل ملف CSV."
      fetchTransactions={fetchStationTransactions}
      exportTransactions={exportStationTransactions}
    />
  );
}
