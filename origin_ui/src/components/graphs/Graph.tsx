/***************************************************************
 *
 * Copyright (C) 2023, Pelican Project, Morgridge Institute for Research
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you
 * may not use this file except in compliance with the License.  You may
 * obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 ***************************************************************/
"use client"

import {useEffect, useState} from "react";
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    TimeScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    ChartOptions,
    Colors
} from 'chart.js';


import zoomPlugin from 'chartjs-plugin-zoom';
import 'chartjs-adapter-luxon';

import {BoxProps} from "@mui/material";

import {Line} from "react-chartjs-2";
import {Box, Skeleton, Typography} from "@mui/material";

import {getDataFunction} from "@/components/graphs/prometheus";
import {ChartData} from "chart.js";

const defaultOptions: Partial<ChartOptions<"line">> = {
    scales: {
        x: {
            type: 'time',
                time: {
                round: 'second',
            }
        }
    }
}

interface GraphProps {
    getData: getDataFunction;
    drawer?: any;
    options?: ChartOptions<"line">
    boxProps?: BoxProps;
}

export default function Graph({getData, options, boxProps, drawer}: GraphProps) {

    let [data, _setData] = useState<ChartData<"line", any, any>>()
    let [loading, setLoading] = useState<boolean>(true)
    let [error, setError] = useState<string>("")


    async function setData() {
        let response = await getData()
        _setData(response)
        setLoading(false)
        if(response.datasets[0].data.length == 0){
            let date = new Date(Date.now()).toLocaleTimeString()
            setError(`No data returned by database as of ${date}; plot will auto-refresh`)
        } else {
            setError("")
        }
    }

    useEffect(() => {

        ChartJS.register(
            CategoryScale,
            LinearScale,
            PointElement,
            LineElement,
            Title,
            Tooltip,
            Legend,
            TimeScale,
            zoomPlugin,
            Colors
        );


        // Do the initial data fetch
        setData()

        // Refetch the data every minute
        const interval = setInterval(() => setData(), 60000);
        return () => clearInterval(interval);

    }, [getData])


    if(loading || !data){
        return <Skeleton sx={{borderRadius: "1"}} variant={"rectangular"} width={"100%"} height={"100%"} />
    }

    return (
        <Box>
            <Box m={"auto"} {...boxProps}>
                <Line
                    data={data}
                    options={{
                        ...defaultOptions,
                        ...options
                    }}
                />
            </Box>
            <Box display={"flex"}>
                {drawer ? drawer : <></>}
            </Box>
            <Box display={"flex"}>
                <Typography m={"auto"} color={"red"} variant={"body2"}>{error}</Typography>
            </Box>
        </Box>
    )

}
