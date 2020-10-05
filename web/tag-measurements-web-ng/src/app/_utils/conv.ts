import * as moment from "moment";

export function vh(v) {
    const h = Math.max(document.documentElement.clientHeight, window.innerHeight || 0);
    return (v * h) / 100;
}

export function vw(v) {
    const w = Math.max(document.documentElement.clientWidth, window.innerWidth || 0);
    return (v * w) / 100;
}

export function downloadCSVFile(data: any[], filename: string) {
    let csv = '';

    for (const tagDataRow of data) {
        csv = csv.concat(`Date,Time,${tagDataRow.name},`);
    }
    csv = csv.concat('\n');
    for (let i = 0; i < data[0].x.length; i++) {
        for (const tagDataRow of data) {
            csv = csv.concat(`${moment(tagDataRow.x[i]).format('YYYY-MM-DD').toString()},${moment(tagDataRow.x[i]).format('HH:mm:ss').toString()},${tagDataRow.y[i].toFixed(6)},`);
        }
        csv = csv.concat('\n');
    }

    const a = document.createElement('a');
    const blob = new Blob([csv], {type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);

    a.href = url;
    a.download = filename;
    a.click();
    window.URL.revokeObjectURL(url);
    a.remove();
}
