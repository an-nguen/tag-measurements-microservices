import {throwError} from "rxjs";
import * as moment from "moment";

interface Point {
    x: number,
    y: number
}

export function approximateTagData(tagData: any, epsilon: number) {
    const points: Point[] = [];
    for (let i = 0; i < tagData.x.length; i++) {
        const p: Point = {
            x: moment(tagData.x[i]).unix(),
            y: tagData.y[i]
        };
        points.push(p);
    }
    let res = douglasPeucker(points, epsilon);
    tagData.x.splice(0, tagData.x.length-1);
    tagData.y.splice(0, tagData.y.length-1);
    for (let p of res) {
        tagData.x.push(p.x);
        tagData.y.push(p.y);
    }
    return tagData;
}

function douglasPeucker(p: Point[], epsilon: number): Point[] {
    let dmax = 0.0; let index = 0;
    let end = p.length - 1;
    for (let i = 0; i < end; i++) {
        let d = perpendicularDistance(p[i], p[0], p[end])
        if (d > dmax) {
            index = i;
            dmax = d;
        }
    }

    let res = [];
    if (dmax > epsilon) {
        let recResults1 = douglasPeucker(p.slice(index+1), epsilon);
        let recResults2 = douglasPeucker(p.slice(index, end), epsilon);

        res.push(...recResults1.slice(0, recResults1.length - 1));
        res.push(...recResults2);
        if (res.length < 2) {
            console.error("Problem assembling output");
        }
    } else {
        res.splice(0, res.length-1)
        res.push(p[0], p[end])
    }
    return res
}

function perpendicularDistance(pt: Point, lineStart: Point, lineEnd: Point): number {
    let dx = lineEnd.x - lineStart.x
    let dy = lineEnd.y - lineStart.y

    let mag = Math.pow(Math.pow(dx, 2.0) + Math.pow(dy, 2.0), 0.5)
    if (mag > 0.0) {
        dx /= mag;
        dy /= mag;
    }

    let pvx = pt.x - lineStart.x;
    let pvy = pt.y - lineStart.y;

    let pvdot = dx*pvx + dy*pvy;

    let ax = pvx - pvdot*dx;
    let ay = pvy - pvdot*dy;

    return Math.pow(Math.pow(ax, 2.0) + Math.pow(ay, 2.0), 0.5);
}